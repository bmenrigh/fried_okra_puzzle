#!/usr/bin/perl

use strict;
use warnings;

# Turn on stdout flushing so we can use tee nicely
$| = 1;

# Goal will always be 0,0

my $LIMIT = 4; # max sum of adjacent diameters

my %dia = (
    'R' => 1,
    'B' => 2,
    'G' => 3,
    'Y' => 3
    );

my %color_names = (
    'R' => 'Red',
    'B' => 'Blue',
    'G' => 'Green',
    'Y' => 'Yellow'
    );

my %opp_dir = (
    'up' => 'down',
    'down' => 'up',
    'left' => 'right',
    'right' => 'left'
    );

my %states = (); # All states found while solving


# Carl's posted puzzles
my $puz_1 = ' Y RGR    B B B ';
my $puz_2 = ' B GRBR  RRGY   ';
my $puz_3 = ' Y   R RG   R   ';
my $puz_4 = '      R YR     R';
my $puz_5 = '   RR G G            Y  R'; # 5x5
my $puz_6 = '   Y R  RR       R  R    '; # 5x5

my $hardest_1 = 'BB  B  R R  G Y '; # 4x4 (R:2, G:1, B:3) 77 moves
my $hardest_2 = 'BR G YR  R  BRG '; # 4x4 (R:4, G:2, B:2) 54 moves
my $hardest_3 = '   G G   R    Y      R  R'; # 5x5 (R:3, G:2, B:0) 243 moves
my $hardest_4 = 'R R   R  R     Y'; # 4x4 (R:4, G:0, B:0) 19 moves
my $hardest_5 = 'RR     RR   Y R '; # 4x4 (R:5, G:0, B:0) 24 moves
my $hardest_6 = 'R R  RRR  R   Y '; # 4x4 (R:6, G:0, B:0) 23 moves
my $hardest_7 = 'RR R     Y RR RR'; # 4x4 (R:7, G:0, B:0) 16 moves
my $hardest_8 = 'R G     GR    YR'; # 4x4 (R:3, G:2, B:0) 53 moves
my $hardest_9 = ' R R Y  G  G RG '; # 4x4 (R:3, G:3, B:0) 52 moves
my $hardest_10 = ' GRG  Y RBRBG  B'; # 4x4 (R:3, G:3, B:3) 41 moves
my $hardest_11 = 'G  G B  RBRYB   '; # 4x4 (R:2, G:2, B:3) 38 moves
my $hardest_12 = ' Y  B BBR B G R '; # 4x4 (R:2, G:1, B:4) 84 moves
my $hardest_13 = ' Y  G BB  B BBRR'; # 4x4 (R:2, G:1, B:5) 78 moves


my $d; # Used later

#$d = find_solution($ring_1, 9, 9);
#print_solution_moves($d, 9, 9);
#exit(0);


# my @solved;
# my %pc = (
#     'R' => 4,
#     'G' => 2,
#     'B' => 2
#     );

# enumerate_solved(\%pc, 4, 4, \@solved);
# $d = explore_from_solved(\@solved, 4, 4);

# foreach my $b (grep {$states{$_} == $d} keys %states) {
#     warn '====== Example Hardest Board =====', "\n", board_2d(\$b, 4, 4), "\n";
#     warn sprintf('\'%s\'; # %dx%d (R:%d, G:%d, B:%d) %d moves',
#                  $b, 4, 4, $pc{'R'}, $pc{'G'}, $pc{'B'}, $d), "\n";
# }
# exit(0);

for (my $tpc = 1; $tpc < 15; $tpc++) {
    for (my $rc = 1; $rc < 16; $rc++) {
        for (my $gc = 0; $gc < 16; $gc++) {
            for (my $bc = 0; $bc < 16; $bc++) {

                next if ($bc == 1);
                next if ($rc + $gc + $bc != $tpc);

                my @solved;
                my %pc = (
                    'R' => $rc,
                    'G' => $gc,
                    'B' => $bc
                    );

                enumerate_solved(\%pc, 4, 4, \@solved);
                #foreach my $b (@solved) {
                #    warn '====== Solved Board =====', "\n", board_2d(\$b, 4, 4), "\n";
                #}

                # Count how many solved states can be scrambled
                my $scrambleable_solved = 0;
                foreach my $b (@solved) {
                    if ((can_move(\$b, 4, 4, 0, 0, 'down') == 1) ||
                        (can_move(\$b, 4, 4, 0, 0, 'right') == 1)) {
                        $scrambleable_solved++;
                    }
                }

                my $fcount = 0;
                $d = explore_from_solved(\@solved, 4, 4, \$fcount);

                my @furthest = grep {$states{$_} == $d} keys %states;
                my $sym_b = 0;
                foreach my $b (sort @furthest) {
                    if ($b eq transpose(\$b, 4, 4)) {
                        $sym_b++;
                    }
                }

                print sprintf('# %dx%d with %d pieces (R:%d, G:%d, B:%d) (solved:%d, scrambleable:%d, found:%d) has %d (symmetric:%d) furthest requiring %d moves',
                              4, 4, $tpc + 1, $pc{'R'}, $pc{'G'}, $pc{'B'},
                              (scalar @solved), $scrambleable_solved, $fcount, (scalar @furthest), $sym_b, $d), "\n";
                foreach my $b (sort @furthest) {
                    print '====== Example Furthest Scramble =====', "\n";
                    print board_2d(\$b, 4, 4), "\n";
                    last;
                }
                print "\n";
            }
        }
    }
}


sub print_solution_moves {
    my $d = shift;
    my $N = shift;
    my $M = shift;

    return if ($d <= 0);

    my @boards = ();
    my @steps = ();

    my $b;
    foreach my $cb (grep {$states{$_} == $d} keys %states) {
        if (get(\$cb, $N, $M, 0, 0) eq 'Y') {
            $b = $cb;
        }
    }

    push @boards, $b;
    while ($d > 0) {
        for (my $x = 0; $x < $N; $x++) {
            for (my $y = 0; $y < $M; $y++) {
                foreach my $dir ('up', 'down', 'left', 'right') {
                    if (can_move(\$b, $N, $M, $x, $y, $dir) == 1) {
                        my $nb = $b;
                        move(\$nb, $N, $M, $x, $y, $dir);

                        next unless (exists $states{$nb});
                        next if ($states{$nb} >= $d);

                        my ($nx, $ny) = find_move_target(\$b, $N, $M, $x, $y, $dir);
                        my $l = get(\$b, $N, $M, $x, $y);
                        my $step = sprintf('Step %d: Move %s at (%d,%d) %s to (%d,%d)',
                                           $d, $color_names{$l}, $nx + 1, $ny + 1, $opp_dir{$dir}, $x + 1, $y + 1);

                        push @steps, $step;
                        push @boards, $nb;
                        $b = $nb;

                        $d = $states{$b};
                    }
                }
            }
        }
    }

    push @steps, 'Step 0: Start with this board';

    for (my $i = (scalar @boards) - 1; $i >= 0; $i--) {
        print $steps[$i], "\n", board_2d(\$boards[$i], $N, $M), "\n\n";
    }
}



sub find_solution {
    my $board = shift;
    my $N = shift;
    my $M = shift;

    %states = ();

    $states{$board} = 0;

    my $d = 0; # depth
    my $done = 0;
    my $f;
    while ($done == 0) {
        $f = 0; # moves found
        #warn '====== At depth ', $d, ' =====', "\n";
        foreach my $b (grep {$states{$_} == $d} keys %states) {
            #warn '====== Working on =====', "\n", board_2d(\$b), "\n";
            for (my $x = 0; $x < $N; $x++) {
                for (my $y = 0; $y < $M; $y++) {
                    foreach my $dir ('up', 'down', 'left', 'right') {
                        if (can_move(\$b, $N, $M, $x, $y, $dir) == 1) {
                            my $nb = $b;
                            move(\$nb, $N, $M, $x, $y, $dir);

                            next if (exists $states{$nb});
                            $states{$nb} = $d + 1;

                            $f++;

                            if (get(\$nb, $N, $M, 0, 0) eq 'Y') {
                                warn 'Found solution in ', ($d + 1), ' moves!', "\n";
                                $done = 1;
                                return $d + 1;
                            }
                        }
                    }
                }
            }
        }

        if ($f > 0) {
            $d++;
        } else {
            warn 'Got to depth ', $d, "\n";
            $done = 1;
        }
    }

    return -1;
}


sub explore_from_solved {
    my $solref = shift;
    my $N = shift;
    my $M = shift;
    my $fcountref = shift;

    %states = ();

    $$fcountref = 0;

    my $sc = 0;
    foreach my $b (@{$solref}) {
        $states{$b} = 0;

        $sc++;
    }
    #warn 'Exploring from ', $sc, ' starting states', "\n";

    my $d = 0; # depth
    my $done = 0;
    my $f;
    while ($done == 0) {
        $f = 0; # moves found
        #warn '====== At depth ', $d, ' =====', "\n";
        foreach my $b (grep {$states{$_} == $d} keys %states) {
            #warn '====== Working on =====', "\n", board_2d(\$b), "\n";
            for (my $x = 0; $x < $N; $x++) {
                next if (($d == 0) && ($x != 0));

                for (my $y = 0; $y < $M; $y++) {
                    next if (($d == 0) && ($y != 0));

                    foreach my $dir ('up', 'down', 'left', 'right') {
                        next if (($d == 0) && (($dir ne 'down') && ($dir ne 'right')));

                        if (can_move(\$b, $N, $M, $x, $y, $dir) == 1) {
                            my $nb = $b;
                            move(\$nb, $N, $M, $x, $y, $dir);

                            next if (exists $states{$nb});
                            $states{$nb} = $d + 1;

                            $f++;

                        }
                    }
                }
            }
        }

        if ($f > 0) {
            #warn 'Depth ', $d, ': found ', $f, ' new states', "\n";
            $d++;

            $$fcountref += $f;
        } else {
            #warn 'Got to depth ', $d, "\n";
            $done = 1;

            return $d;
        }
    }

    return -1;
}


sub set {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $x = shift;
    my $y = shift;
    my $l = shift;

    if ((($x < 0) || ($x >= $N)) ||
        (($y < 0) || ($y >= $M))) {
        return undef;
    }

    substr($$boardref, ($y * $N) + $x, 1) = $l;
}


sub get {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $x = shift;
    my $y = shift;

    if ((($x < 0) || ($x >= $N)) ||
        (($y < 0) || ($y >= $M))) {
        return undef;
    }

    return substr($$boardref, ($y * $N) + $x, 1);
}


sub get_dia {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $x = shift;
    my $y = shift;

    if ((($x < 0) || ($x >= $N)) ||
        (($y < 0) || ($y >= $M))) {
        return 0; # Could also make the outside 'take up space'
    }

    my $t = get($boardref, $N, $M, $x, $y);

    return 0 unless (defined $t);
    return 0 if ($t eq ' ');

    return $dia{$t};
}


sub will_fit {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $x = shift;
    my $y = shift;
    my $l = shift;

    my $t = get($boardref, $N, $M, $x, $y);

    #warn sprintf('will_fit("%s", %d, %d, "%s")', $$boardref, $x, $y, $l), "\n";
    #warn '======', "\n", board_2d($boardref), "\n", '======', "\n";

    return 0 unless (defined $t);
    return 0 unless ($t eq ' ');

    my $r = $dia{$l};

    unless (defined $r) {
        die 'Unknown diameter for "', $l, '"', "\n";
    }

    # Check up
    return 0 if ($r + get_dia($boardref, $N, $M, $x, $y - 1) > $LIMIT);

    # Check down
    return 0 if ($r + get_dia($boardref, $N, $M, $x, $y + 1) > $LIMIT);

    # Check left
    return 0 if ($r + get_dia($boardref, $N, $M, $x - 1, $y) > $LIMIT);

    # Check right
    return 0 if ($r + get_dia($boardref, $N, $M, $x + 1, $y) > $LIMIT);

    return 1;
}


sub can_move {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $x = shift;
    my $y = shift;
    my $dir = shift;

    my $l = get($boardref, $N, $M, $x, $y);
    return 0 unless (defined $l);
    return 0 if ($l eq ' ');

    my ($dx, $dy);
    if ($dir eq 'up') {
        ($dx, $dy) = (0, -1);
    } elsif ($dir eq 'down') {
        ($dx, $dy) = (0, 1);
    } elsif ($dir eq 'left') {
        ($dx, $dy) = (-1, 0);
    } elsif ($dir eq 'right') {
        ($dx, $dy) = (1, 0);
    } else {
        die 'Unknown direction: ', $dir, "\n";
    }

    my ($cx, $cy) = ($x, $y);
    my $done = 0;
    my $found = 0;
    my $good = 0;
    do {
        ($cx, $cy) = ($cx + $dx, $cy + $dy);

        my $t = get($boardref, $N, $M, $cx, $cy);

        if (defined $t) {
            if ($t eq ' ') {
                $done = 1;
                $good = 1;
            } else {
                $found++;
            }
        } else {
            $done = 1;
        }
    }  while ($done == 0);

    if (($found > 0) && ($good == 1)) {
        return will_fit($boardref, $N, $M, $cx, $cy, $l);
    } else {
        return 0;
    }

    die 'Got to end of can_move()', "\n";
}


sub move {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $x = shift;
    my $y = shift;
    my $dir = shift;

    my $l = get($boardref, $N, $M, $x, $y);
    die 'Invalid location!', "\n" unless (defined $l);
    die 'No piece at that location!', "\n", if ($l eq ' ');

    die 'Unable to move!', "\n" unless (can_move($boardref, $N, $M, $x, $y, $dir) == 1);


    my ($nx, $ny) = find_move_target($boardref, $N, $M, $x, $y, $dir);
    # Clear this spot
    set($boardref, $N, $M, $x, $y, ' ');
    set($boardref, $N, $M, $nx, $ny, $l);
}


sub find_move_target {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $x = shift;
    my $y = shift;
    my $dir = shift;

    my $l = get($boardref, $N, $M, $x, $y);
    die 'Invalid location!', "\n" unless (defined $l);
    die 'Can not find target: No piece at that location!', "\n", if ($l eq ' ');

    die 'Unable to move!', "\n" unless (can_move($boardref, $N, $M, $x, $y, $dir) == 1);

    my ($dx, $dy);
    if ($dir eq 'up') {
        ($dx, $dy) = (0, -1);
    } elsif ($dir eq 'down') {
        ($dx, $dy) = (0, 1);
    } elsif ($dir eq 'left') {
        ($dx, $dy) = (-1, 0);
    } elsif ($dir eq 'right') {
        ($dx, $dy) = (1, 0);
    } else {
        die 'Unknown direction: ', $dir, "\n";
    }

    my ($cx, $cy) = ($x, $y);
    my $done = 0;
    do {
        ($cx, $cy) = ($cx + $dx, $cy + $dy);

        my $t = get($boardref, $N, $M, $cx, $cy);

        if (defined $t) {
            if ($t eq ' ') {
                $done = 1;
                return ($cx, $cy);
            }
        } else {
            die 'Tried to move off end of board!', "\n";
        }
    }  while ($done == 0);
}


sub board_2d {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;

    my @rows = ();
    for (my $y = 0; $y < $M; $y++) {
        push @rows, substr($$boardref, $y * $N, $N);
    }

    my $bstr = join("\n", @rows);

    $bstr =~ s/ /./g;

    return $bstr;
}


sub enumerate_solved {
    my $pcref = shift;
    my $N = shift;
    my $M = shift;
    my $solref = shift;

    my $blank = (' ' x ($N * $M));

    my $start = $blank;
    set(\$start, $N, $M, 0, 0, 'Y');

    rec_enum(\$start, $N, $M, $pcref, 1, $solref);
}


sub rec_enum {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;
    my $pcref = shift;
    my $o = shift;
    my $solref = shift;

    my ($cx, $cy) = ($o % $N, int($o / $N));

    # If the board is filled out
    if ($cy >= $M) {
        my $leftover = 0;

        foreach my $l (keys %{$pcref}) {
            if ($pcref->{$l} > 0) {
                $leftover = 1;
                last;
            }
        }

        if ($leftover == 0) {
            my $b = $$boardref;
            push @{$solref}, $b;
        }

        return;
    }

    foreach my $l (keys %{$pcref}) {
        if ($pcref->{$l} > 0) {
            if (will_fit($boardref, $N, $M, $cx, $cy, $l) == 1) {
                my $b = $$boardref;
                my %pc = %{$pcref};

                set(\$b, $N, $M, $cx, $cy, $l);
                $pc{$l} -= 1;

                # Put this piece in
                rec_enum(\$b, $N, $M, \%pc, $o + 1, $solref);
            }
        }
    }
    # Leave this spot blank
    rec_enum($boardref, $N, $M, $pcref, $o + 1, $solref);
}


sub transpose {
    my $boardref = shift;
    my $N = shift;
    my $M = shift;

    my $nb = '';
    for (my $x = 0; $x < $N; $x++) {
        for (my $y = 0; $y < $M; $y++) {
            $nb .= substr($$boardref, ($y * $N) + $x, 1);
        }
    }

    return $nb;
}
