<?php

namespace NodaSoft\Operation\Command;

use NodaSoft\GenericDto\Factory\ComplaintStatusChangedMessageContentListFactory;
use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\InitialData\NotifyComplaintStatusChangedInitialData;
use NodaSoft\Operation\Result\Result;
use NodaSoft\Operation\Result\NotifyComplaintStatusChangedResult;

class NotifyComplaintStatusChangedCommand implements Command
{
    /** @var Messenger */
    private $email;

    /** @var Messenger */
    private $sms;

    public function setEmail(Messenger $email): void
    {
        $this->email = $email;
    }

    public function setSms(Messenger $sms): void
    {
        $this->sms = $sms;
    }

    /**
     * @param NotifyComplaintStatusChangedInitialData $data
     * @return NotifyComplaintStatusChangedResult
     * @throws \Exception Previous status required, 500
     */
    public function execute(InitialData $data): Result
    {
        $result = new NotifyComplaintStatusChangedResult();
        $complaint = $data->getComplaint();
        $client = $complaint->getClient();
        $reseller = $complaint->getReseller();

        $contentFactory = new ComplaintStatusChangedMessageContentListFactory();

        try {
            $contentList = $contentFactory->composeContentList($complaint);
        } catch (\Exception $e) {
            throw new \Exception($e->getMessage(), 500, $e);
        }

        $message = new Message($data->getNotification(), $contentList);

        foreach ($reseller->getEmployees() as $employee) {
            $result->addEmployeeEmailResult($this->email->send($message, $employee, $reseller));
        }

        $result->setClientEmailResult($this->email->send($message, $client, $reseller));
        $result->setClientSmsResult($this->sms->send($message, $client, $reseller));

        return $result;
    }
}
