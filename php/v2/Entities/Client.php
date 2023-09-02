<?php

namespace NW\WebService\References\Operations\Notification\Entities;


/**
 * @property Seller $seller
 * @property string $email
 * @property mixed $mobile // Не понятно какой тип должен быть
 */
class Client extends Contractor
{
    protected $email;
    protected $mobile;

    /**
     * @return mixed
     */
    public function getEmail()
    {
        return $this->email;
    }

    /**
     * @return mixed
     */
    public function getMobile()
    {
        return $this->mobile;
    }




}