<?php

require_once 'src/Stubs/functions.php';
require_once '../../../php/v2/others.php';
require_once 'vendor/autoload.php';

use Israil\V2\ReturnOperation;

echo print_r((new ReturnOperation)->doOperation(), true);
