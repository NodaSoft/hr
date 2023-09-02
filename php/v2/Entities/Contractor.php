<?php

namespace  NW\WebService\References\Operations\Notification\Entities;
use EntityNotFoundException;

abstract class Contractor
{
    const TYPE_CUSTOMER = 0;
    protected $id;
    protected $type;
    protected $name;

    public function __construct(int $resellerId)
    {
        $this->id = $resellerId;
    }

    /**
     * @param int $id
     * @return Contractor|null
     */
    public static function getById(int $id): ?Contractor
    {
        try {
            // К примеру ищем по базе и не находим, выдаём null, т.к. в коде есть проверки.
            $entity = new static($id);
        } catch (EntityNotFoundException $exception) {
            return null;
        }

        return $entity;
    }

    /**
     * @return string
     */
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

    /**
     * @return mixed
     */
    public function getType()
    {
        return $this->type;
    }

    /**
     * @return mixed
     */
    public function getName()
    {
        return $this->name;
    }

    public function getId(): int
    {
        return $this->id;
    }




}