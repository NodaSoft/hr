<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

/**
 * @property-read Seller $seller
 * @property-read int $type
 */
class Contractor extends AbstractContractor
{
    const TYPE_CUSTOMER = 0;

    private $_type;
    private $_seller;

    /**
     * @param int $id
     * @throws Exception
     */
    public function __construct(int $id)
    {
        parent::__construct($id);
        $this->_type = 0;

        $sellerId = 123; // fake id
        if (!$this->_seller = new Seller($sellerId)) {
            throw new Exception("Seller with id = '{$sellerId}' not found");
        }
    }

    /**
     * @param $name
     * @throws Exception
     */
    public function __get($name)
    {
        switch ($name) {
            case 'type':
                return $this->_type;
            case 'Seller':
                return $this->_seller;
            default:
                return parent::__get($name);
        }
    }
}
