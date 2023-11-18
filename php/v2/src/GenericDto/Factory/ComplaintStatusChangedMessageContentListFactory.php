<?php

namespace NodaSoft\GenericDto\Factory;

use NodaSoft\DataMapper\Entity\Complaint;
use NodaSoft\GenericDto\Dto\ComplaintStatusChangedMessageContentList;

class ComplaintStatusChangedMessageContentListFactory
{
    public function composeContentList(
        Complaint $complaint
    ): ComplaintStatusChangedMessageContentList {
        $client = $complaint->getClient();
        $consumption = $client->getConsumption();
        $creator = $complaint->getCreator();
        $expert = $complaint->getExpert();
        $currentStatus = $complaint->getStatus();
        $previousStatus = $complaint->getPreviousStatus();
        $contentList = new ComplaintStatusChangedMessageContentList();

        $contentList->setComplaintId($complaint->getId());
        $contentList->setComplaintNumber($complaint->getNumber());
        $contentList->setCreatorId($creator->getId());
        $contentList->setCreatorName($creator->getFullName());
        $contentList->setExpertId($expert->getId());
        $contentList->setExpertName($expert->getFullName());
        $contentList->setClientId($client->getId());
        $contentList->setClientName($client->getFullName());
        $contentList->setConsumptionId($consumption->getId());
        $contentList->setConsumptionNumber($consumption->getNumber());
        $contentList->setAgreementNumber($consumption->getAgreementNumber());
        $contentList->setDate((new \DateTime())->format(DATE_W3C));
        $contentList->setCurrentStatus($currentStatus->getName());
        $contentList->setPreviousStatus($previousStatus->getName());

        return $contentList;
    }
}
