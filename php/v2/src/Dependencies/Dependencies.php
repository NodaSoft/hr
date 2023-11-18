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
    /** @var ? Request  */
    private $request = null;

    /** @var ? Messenger */
    private $emailService = null;

    /** @var ? Messenger */
    private $smsService = null;

    /** @var ? MapperFactory */
    private $mapperFactory = null;

    public function __construct(
        ? Request $request = null,
        ? Messenger $emailService = null,
        ? Messenger $smsService = null,
        ? MapperFactory $mapperFactory = null
    ) {
        $this->request = $request;
        $this->emailService = $emailService;
        $this->smsService = $smsService;
        $this->mapperFactory = $mapperFactory;
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
