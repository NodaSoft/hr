<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Contracts;

/**
 * Class Contractor.
 *
 */
class Contractor
{
    public const TYPE_CUSTOMER = 0;
    /**
     * @var int
     */
    private int $id;
    /**
     * @var int
     */
    private int $type;
    /**
     * @var string
     */
    private string $name;
    /**
     * @var string
     */
    private string $email;
    /**
     * @var bool
     */
    private bool $mobile;

    /**
     * Contractor constructor.
     *
     * @param int $resellerId
     * @param int $type
     * @param string $name
     * @param string $email
     * @param bool $mobile
     */
    public function __construct(
        int    $resellerId,
        int    $type = self::TYPE_CUSTOMER,
        string $name = 'fake name',
        string $email = 'fake_email@gmail.com',
        bool   $mobile = false
    )
    {
        $this->id = $resellerId;
        $this->type = $type;
        $this->name = $name;
        $this->email = $email;
        $this->mobile = $mobile;
    }

    /**
     * get Contractor obj by reseller id
     *
     * @param int $resellerId
     * @return Contractor
     */
    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
    }

    /**
     * get Contractor full name
     *
     * @return string
     */
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

    /**
     * get Contractor type
     *
     * @return int
     */
    public function getType(): int
    {
        return $this->type;
    }

    /**
     * get Contractor Id
     *
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * get Contractor Email
     *
     * @return string
     */
    public function getEmail(): string
    {
        return $this->email;
    }

    /**
     * get Contractor If
     *
     * @return bool
     */
    public function isMobile(): bool
    {
        return $this->mobile;
    }
}