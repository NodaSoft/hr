<?php

// declare strict type

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */

// не совсем понял зачем в одном файле делать много классов и функций не привязанному к файлу - это к коду всего файла
class Contractor
{
    const TYPE_CUSTOMER = 0;

    //Свойства private или protected // класс должен быть закрыт для изменений
    public $id;
    public $type;
    public $name;

    // Где конструктор? Должна быть жесткая типизация иначе получается смесь говна

    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
}

// Кек, создаются разные классы Seller Employee, но различий в них нету
// final
class Seller extends Contractor
{
}

// final
class Employee extends Contractor
{
}


// Completed, Pending, Rejected - сделать как константы
// final
class Status
{
//    Открпытие всему миру
    public $id, $name;

    public static function getName(int $id): string
    {
        $a = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $a[$id];
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    // Возращаемый тип, я за жесткую типизацию
    public function getRequest($pName)
    {
        //Ахуенно, глобальный переменные, с таким успехом предалгаю вернуться на 15 лет назад
        return $_REQUEST[$pName];
    }
}

// Если так уж надо, то сделать ValueObject, а не функцию не приязаную никчему
// Ввозвращаемыый тип
function getResellerEmailFrom()
{
    return 'contractor@example.com';
}

// Если так уж надо, то сделать ValueObject, а не функцию не приязаную никчему
// Ввозвращаемыый тип
// Указать тип переменных  $resellerId, $event
// Зачем тут поля?
function getEmailsByPermit($resellerId, $event)
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}


//final
class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';

    // Убрать, он не используется
    // Так же лучше создавать отдельный класс на отдельное событие с отдельными полями, Чем все в одном месте хранить
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}