<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

class SellerNotFoundException extends \Exception {
    public function __construct()
    {
        parent::__construct("Seller not found", 400);
    }
}
