<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor{
    
    protected $id;
    
    public $type;
    public $name;
    public $mobile;
    public $email;
    
    public function __construct( $id )
    {
        // @todo should be filtered by type too
        $data = $this->getDataFromSomewhere( $id );
        
        if( empty( $data ) ){
            throw new \Exception( static::class . ' not found!', 400 );
        }
        
        $this->id     = $data['id'] ?? null;
        $this->type   = $data['type'] ?? null;
        $this->name   = $data['name'] ?? null;
        $this->mobile = $data['mobile'] ?? null;
        $this->email  = $data['email'] ?? null;
    }
    
    public static function getById( int $id ): self
    {
        return new static( $id ); // fakes the getById method
    }
    
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
    
    /**
     * Receives data from some storage
     *
     * @param $id
     *
     * @return array
     */
    protected function getDataFromSomewhere( $id ): array
    {
        // stub
        
        return [];
    }
    
    /**
     * @return mixed
     */
    public function getId()
    {
        return $this->id;
    }
    
    /**
     * @return mixed|null
     */
    public function getType()
    {
        return $this->type;
    }
    
    /**
     * @return mixed|null
     */
    public function getName()
    {
        return $this->name;
    }
    
    /**
     * @return mixed|null
     */
    public function getMobile()
    {
        return $this->mobile;
    }
    
    /**
     * @return mixed|null
     */
    public function getEmail()
    {
        return $this->email;
    }
}
