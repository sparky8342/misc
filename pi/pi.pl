#!/usr/bin/perl
use strict;
use warnings;
use Math::BigFloat;

# Gauss-Legendre Algorithm for calculating pi

use constant DIGITS => 1000;

Math::BigFloat->accuracy( DIGITS + 2 );

my $two = Math::BigFloat->new(2);

my $a      = Math::BigFloat->new(1);
my $b      = Math::BigFloat->new(1)->bdiv( Math::BigFloat->new(2)->bsqrt() );
my $t      = Math::BigFloat->new(1)->bdiv( Math::BigFloat->new(4) );
my $x      = Math::BigFloat->new(1);
my $last_t = Math::BigFloat->new(0);

for ( my $i = 0 ; $i < int( log(DIGITS) / log(2) ) ; $i++ ) {
    my $y = $a->copy();
    $a->badd($b);
    $a->bdiv($two);
    $b->bmul($y);
    $b->bsqrt();
    $y->bsub($a);
    $y->bpow(2);
    $y->bmul($x);
    $t->bsub($y);
    $x->bmul($two);
}

$a->badd($b);
$a->bpow(2);
my $bot = Math::BigFloat->new(4);
$bot->bmul($t);
$a->bdiv($bot);
$a->bfround( -DIGITS );
print "$a\n";
