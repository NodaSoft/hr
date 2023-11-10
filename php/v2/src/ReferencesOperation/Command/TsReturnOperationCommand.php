<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\Message\Message;
use NodaSoft\Message\Messenger;
use NodaSoft\Message\Result;
use NodaSoft\Message\Template\ComplaintStatus;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\Message\Factory\MessageFactory;
use NodaSoft\ReferencesOperation\InitialData\TsReturnInitialData;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\TsReturnOperationResult;

class TsReturnOperationCommand implements ReferencesOperationCommand
{
    /** @var TsReturnOperationResult */
    private $result;

    /** @var InitialData */
    private $initialData;

    /** @var Messenger */
    private $mail;

    /** @var Messenger */
    private $sms;

    /**
     * @param TsReturnOperationResult $result
     * @return void
     */
    public function setResult(ReferencesOperationResult $result): void
    {
        $this->result = $result;
    }

    public function setInitialData(InitialData $initialData): void
    {
        $this->initialData = $initialData;
    }

    public function setMail(Messenger $mail): void
    {
        $this->mail = $mail;
    }

    public function setSms(Messenger $sms): void
    {
        $this->sms = $sms;
    }

    /**
     * @return TsReturnOperationResult
     */
    public function execute(): ReferencesOperationResult
    {
        /** @var TsReturnInitialData $initialData */
        $initialData = $this->initialData;
        $client = $initialData->getClient();

        $complaintTemplate = new MessageFactory(new ComplaintStatus());
        $reseller = $initialData->getReseller();

        foreach ($initialData->getEmployees() as $employee) {
            $message = $complaintTemplate->makeMessage($employee, $reseller, $initialData);
            $result = $this->mail->send($message);
            $this->result->addEmployeeEmailResult($result);
        }

        $message = $complaintTemplate->makeMessage($client, $reseller, $initialData);
        $result = $this->mail->send($message);
        $this->result->setClientEmailResult($result);

        $result = $this->sendClientSms($initialData, $message);
        $this->result->setClientSmsResult($result);

        return $this->result;
    }

    public function sendClientSms(InitialData $data, Message $message): Result
    {
        if ($data->getNotification()->getName() !== "complaint status changed"
            || is_null($data->getDifferencesTo())) { //todo: replace the condition with a unit logic
            $result = new Result($message->getRecipient(), get_class($this->sms));
            $result->setIsSent(false);
            $result->setErrorMessage("Wrong parameters. The message was not sent.");
            return $result;
        }

        return $this->sms->send($message);
    }
}
