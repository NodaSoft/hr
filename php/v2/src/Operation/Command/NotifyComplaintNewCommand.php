<?php

namespace NodaSoft\Operation\Command;

use NodaSoft\GenericDto\Factory\ComplaintNewMessageContentListFactory;
use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\InitialData\NotifyComplaintNewInitialData;
use NodaSoft\Operation\Result\Result;
use NodaSoft\Operation\Result\NotifyComplaintNewResult;

class NotifyComplaintNewCommand implements Command
{
    /** @var Messenger */
    private $email;

    public function setEmail(Messenger $email): void
    {
        $this->email = $email;
    }

    /**
     * @param NotifyComplaintNewInitialData $data
     * @return NotifyComplaintNewResult
     */
    public function execute(InitialData $data): Result
    {
        $result = new NotifyComplaintNewResult();
        $complaint = $data->getComplaint();
        $reseller = $complaint->getReseller();

        $contentFactory = new ComplaintNewMessageContentListFactory();
        $contentList = $contentFactory->composeContentList($complaint);

        $message = new Message($data->getNotification(), $contentList);

        foreach ($reseller->getEmployees() as $employee) {
            $result->addEmployeeEmailResult(
                $this->email->send($message, $employee, $reseller)
            );
        }

        return $result;
    }
}
