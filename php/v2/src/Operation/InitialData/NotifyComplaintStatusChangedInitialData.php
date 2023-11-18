<?php

namespace NodaSoft\Operation\InitialData;

use NodaSoft\DataMapper\Entity\Complaint;
use NodaSoft\DataMapper\Entity\Notification;

class NotifyComplaintStatusChangedInitialData implements InitialData
{
    /** @var Notification */
    private $notification;

    /** @var Complaint */
    private $complaint;

    public function getNotification(): Notification
    {
        return $this->notification;
    }

    public function setNotification(Notification $notification): void
    {
        $this->notification = $notification;
    }

    public function getComplaint(): Complaint
    {
        return $this->complaint;
    }

    public function setComplaint(Complaint $complaint): void
    {
        $this->complaint = $complaint;
    }
}
