<?php

    namespace NW\WebService\References\Operations\Notification\BusinessLayer;

    class SendMessageRequest
    {
        public int $resellerId;
        public int $notificationType;
        public int $clientId;
        public int $creatorId;
        public int $expertId;
        public int $complaintId;
        public int $consumptionId;
        public string $complaintNumber;
        public string $consumptionNumber;
        public string $agreementNumber;
        public string $date;
        public DifferenceMessageRequest $differences;

        public function __construct(array $raw) {
            $this->resellerId = (int)$raw['resellerId'];
            $this->notificationType = (int)$raw['notificationType'];
            $this->clientId = (int)$raw['clientId'];
            $this->creatorId = (int)$raw['creatorId'];
            $this->expertId = (int)$raw['expertId'];
            $this->complaintId = (int)$raw['complaintId'];
            $this->consumptionId = (int)$raw['consumptionId'];
            $this->complaintNumber = (string)$raw['complaintNumber'];
            $this->consumptionNumber = (string)$raw['consumptionNumber'];
            $this->agreementNumber = (string)$raw['agreementNumber'];
            $this->date = (string)$raw['date'];
            $this->differences = new DifferenceMessageRequest($raw);
        }

    }
