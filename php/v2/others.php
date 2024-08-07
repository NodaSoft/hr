<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;
    public $id;
    public $type;
    public $name;

    public static function getById(int $resellerId): ?self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
}

class Seller extends Contractor
{
}

class Employee extends Contractor
{
}

class Status
{
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

    public function getRequest($pName)
    {
        return $_REQUEST[$pName];
    }
}

function getResellerEmailFrom(): string
{
    return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}

final class ReturnOperationDTO
{
    private int $resellerId;
    private int $clientId;
    private int $creatorId;
    private int $expertId;
    private int $differencesFrom;
    private int $differencesTo;
    private int $notificationType;
    private int $complaintId;
    private int $consumptionId;
    private int $complaintNumber;
    private string $agreementNumber;
    private string $date;

    public function __construct(
        int $resellerId,
        int $clientId,
        int $creatorId,
        int $expertId,
        int $differencesFrom,
        int $differencesTo,
        int $notificationType,
        int $complaintId,
        int $consumptionId,
        string $complaintNumber,
        string $consumptionNumber,
        string $agreementNumber,
        string $date
    )
    {
        $this->resellerId = $resellerId;
        $this->clientId = $clientId;
        $this->creatorId = $creatorId;
        $this->expertId = $expertId;
        $this->differencesFrom = $differencesFrom;
        $this->differencesTo = $differencesTo;
        $this->notificationType = $notificationType;
        $this->complaintId = $complaintId;
        $this->consumptionId = $consumptionId;
        $this->complaintNumber = $complaintNumber;
        $this->consumptionNumber = $consumptionNumber;
        $this->agreementNumber = $agreementNumber;
        $this->date = $date;
    }

    public function getResellerId(): int
    {
        return $this->resellerId;
    }

    public function getClientId(): int
    {
        return $this->clientId;
    }

    public function getCreatorId(): int
    {
        return $this->creatorId;
    }

    public function getExpertId(): int
    {
        return $this->expertId;
    }

    public function getDifferencesFrom(): int
    {
        return $this->differencesFrom;
    }

    public function getDifferencesTo(): int
    {
        return $this->differencesTo;
    }

    public function getNotificationType(): int
    {
        return $this->notificationType;
    }

    public function getComplaintId(): int
    {
        return $this->complaintId;
    }

    public function getConsumptionId(): int
    {
        return $this->consumptionId;
    }

    public function getComplaintNumber(): int
    {
        return $this->complaintNumber;
    }

    public function getAgreementNumber(): string
    {
        return $this->agreementNumber;
    }

    public function getDate(): string
    {
        return $this->date;
    }
}

class ValidateRequestDataException extends \InvalidArgumentException { }
class NotFoundEntityException extends \DomainException { }
class TemplateException extends \Exception { }