<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

/**
 * @property-read int $id
 * @property-read string $name
 */

abstract class AbstractContractor
{
    protected $_id;
    protected $_name;

    /**
     * @param int $id
     * @throws
     */
    public function __construct(int $id)
    {
        // Если невозможно создать объект с id, то вызываем ошибку
        if ($id == 0) { // для простоты пока укажем здесь якобы неверный id
            throw new ExceptionAPI('Not found', 404);
        }

        $this->_id = $id;
        $this->_name = '';
    }

    /**
     * @return string
     */
    public function getFullName(): string
    {
        return $this->_name . ' ' . $this->_id;
    }

    /**
     * @param $name
     * @return int|null
     * @throws Exception
     */
    public function __get($name)
    {
        switch ($name) {
            case 'id':
                return $this->_id;
            case 'name':
                return $this->_name;
            default:
                throw new Exception("У объект нет свойства {$name}");
        }
    }
}