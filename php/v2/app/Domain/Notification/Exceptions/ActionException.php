<?php
namespace app\Domain\Notification\Exceptions;

class ActionException extends \Exception
{
    /**
     * {@inheritdoc}
     */
    protected $message = 'An error occurred';
}
