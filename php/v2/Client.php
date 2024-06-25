<?php

namespace NW\WebService\References\Operations\Notification;

class Client extends Contractor
{
    const TYPE_CUSTOMER = 0;
    
    // public $mobile;
    
    public function __construct( $id )
    {
        parent::__construct( $id );
        
        // Additional validation
        if( $this->type !== self::TYPE_CUSTOMER ){
            throw new \Exception( 'Client not found!', 400 );
        }
    }
}
