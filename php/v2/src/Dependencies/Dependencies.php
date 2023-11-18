<?php

namespace NodaSoft\Dependencies;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Messenger\Client\EmailClient;
use NodaSoft\Messenger\Client\SmsClient;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Request\HttpRequest;
use NodaSoft\Request\Request;

class Dependencies
{
    /** @var Request  */
    private $request;

    /** @var Messenger */
    private $emailService;

    /** @var Messenger */
    private $smsService;

    /** @var MapperFactory */
    private $mapperFactory;

    public function __construct(
        ? Request $request = null,
        ? Messenger $emailService = null,
        ? Messenger $smsService = null,
        ? MapperFactory $mapperFactory = null
    ) {
        if ($request) $this->request = $request;
        if ($emailService) $this->emailService = $emailService;
        if ($smsService) $this->smsService = $smsService;
        if ($mapperFactory) $this->mapperFactory = $mapperFactory;
    }

    public function getRequest(): Request
    {
        return $this->request ?? new HttpRequest();
    }

    public function getEmailService(): Messenger
    {
        return $this->emailService ?? new Messenger(new EmailClient());
    }

    public function getSmsService(): Messenger
    {
        return $this->smsService ?? new Messenger(new SmsClient());
    }

    public function getMapperFactory(): MapperFactory
    {
        return $this->mapperFactory ?? new MapperFactory();
    }
}
