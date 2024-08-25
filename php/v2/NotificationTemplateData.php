<?php

namespace NW\WebService\References\Operations\Notification;

class NotificationTemplateData
{
    private array $templateData = [
        'COMPLAINT_ID'       => null,
        'COMPLAINT_NUMBER'   => null,
        'CREATOR_ID'         => null,
        'CREATOR_NAME'       => null,
        'EXPERT_ID'          => null,
        'EXPERT_NAME'        => null,
        'CLIENT_ID'          => null,
        'CLIENT_NAME'        => null,
        'CONSUMPTION_ID'     => null,
        'CONSUMPTION_NUMBER' => null,
        'AGREEMENT_NUMBER'   => null,
        'DATE'               => null,
        'DIFFERENCES'        => null,
    ];

    /**
     * @return array
     */
    public function getTemplateData(): array
    {
        return $this->templateData;
    }

    public function validate(): void
    {
        foreach ($this->templateData as $key => $value) {
            if (empty($value)) {
                throw new HttpInternalServerErrorException("Template Data ({$key}) is empty!");
            }
        }
    }


    public function fillWithNotification(NotificationDTO $notification): self
    {
        $this->templateData['COMPLAINT_ID'] = $notification->getComplaintId();
        $this->templateData['COMPLAINT_NUMBER'] = $notification->getComplaintId();
        $this->templateData['CONSUMPTION_ID'] = $notification->getComplaintId();
        $this->templateData['CONSUMPTION_NUMBER'] = $notification->getComplaintId();
        $this->templateData['AGREEMENT_NUMBER'] = $notification->getComplaintId();
        $this->templateData['DATE'] = $notification->getComplaintId();

        return $this;
    }

    public function fillWithCreator(Employee $creatorEmployee)
    {
        $this->templateData['CREATOR_ID']   = $creatorEmployee->getId();
        $this->templateData['CREATOR_NAME'] = $creatorEmployee->getName();
    }

    public function fillWithExpert(Employee $expertEmployee)
    {
        $this->templateData['EXPERT_ID']   = $expertEmployee->getId();
        $this->templateData['EXPERT_NAME'] = $expertEmployee->getName();
    }

    public function fillWithClient(Contractor $client)
    {
        $this->templateData['CLIENT_ID']   = $client->getId();
        $this->templateData['CLIENT_NAME'] = $client->getName();
    }

    public function fillWithDifferences(array $differences)
    {
        $this->templateData['DIFFERENCES'] = $differences;
    }
}