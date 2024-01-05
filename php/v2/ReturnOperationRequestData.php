<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

/**
 * @property-read int $resellerId;
 * @property-read int $notificationType;
 * @property-read int $clientId;
 * @property-read int $creatorId;
 * @property-read int $expertId;
 * @property-read array $differences;
 * @property-read int $complaintId;
 * @property-read string $complaintNumber;
 * @property-read int $consumptionId;
 * @property-read string $consumptionNumber;
 * @property-read string $agreementNumber;
 * @property-read string $date;
 */
class ReturnOperationRequestData
{
    private $types = [
        'resellerId' => ['int', 0],
        'notificationType' => ['int', 0],
        'clientId' => ['int', 0],
        'creatorId' => ['int', 0],
        'expertId' => ['int', 0],
        'differences' => ['array', ['from' => 0, 'to' => 0]],
        'complaintId' => ['int', 0],
        'complaintNumber' => ['string', ''],
        'consumptionId' => ['int', 0],
        'consumptionNumber' => ['string', ''],
        'agreementNumber' => ['string', ''],
        'date' => ['string', '']
    ];

    private $data = [];

    /**
     * @param $data
     * @throws ExceptionAPI
     */
    public function __construct($data)
    {
        if (is_null($data)) {
            // Если вообще не передан параметр data
            throw new ExceptionAPI('Empty request data', 500);
        } elseif (!is_array($data)) {
            // Если параметр data передан не в виде массива
            throw new ExceptionAPI('Request data must be array', 500);
        } else {
            // Устанавливаем свойства объекта, попутно типизируем их
            foreach ($this->types as $field => $type) {
                if (key_exists($field, $data)) {
                    settype($data[$field], $type[0]);
                    if ($field == 'differences') {
                        // Для этого параметра мы ожидаем два элемента массива: from и to
                        if (!key_exists('from', $data[$field]) || !key_exists('to', $data[$field])) {
                            throw new ExceptionAPI("Parameter 'differences' in request data is wrong", 500);
                        } else {
                            // Если в параметре есть лишние данные, то отсекаем их
                            $this->data[$field] = ['from' => $data[$field]['from'], 'to' => $data[$field]['to']];
                        }
                    } else {
                        $this->data[$field] = $data[$field];
                    }
                } else {
                    // Если не передано значение, устанавливаем значение по умолчанию
                    $this->data[$field] = $type[1];
                }
            }
        }
    }

    /**
     * @param $name
     * @return void
     * @throws Exception
     */
    public function __get($name)
    {
        if (!key_exists($name, $this->data)) {
            throw new Exception("У объект нет свойства {$name}");
        }

        return $this->data[$name];
    }
}