<?php
namespace app\Domain\Notification\Models;
class Employee
{
    const TYPE_RESELLER = 1;
    const TYPE_CONTRACTOR = 2;
    const TYPE_CREATOR = 3;
    const TYPE_EXPERT = 4;

    public $id;
    public $type;
    public $name;
    public $email;
    public $mobile;

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
}