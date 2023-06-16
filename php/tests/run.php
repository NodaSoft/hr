<?php

declare(strict_types=1);

use App\Tests\SimpleTest;

require dirname(__DIR__) . '/vendor/autoload.php';

(new SimpleTest())->run();
