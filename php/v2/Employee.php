<?php

namespace NW\WebService\References\Operations\Notification;

class Employee extends Contractor
{
    public function __construct( $id )
    {
        parent::__construct( $id );
        
        // @todo any additional checks here
    }

}
