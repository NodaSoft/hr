<?php

    namespace NW\WebService\References\Operations\Notification;

    use NW\WebService\References\Operations\Notification\BusinessLayer\Enums\ENotificationDifferentType;
    use NW\WebService\References\Operations\Notification\BusinessLayer\NotificationStatusType;
    use NW\WebService\References\Operations\Notification\BusinessLayer\SendMessageRequest;
    use NW\WebService\References\Operations\Notification\BusinessLayer\SendNotificationResponse;
    use NW\WebService\References\Operations\Notification\Controllers\ReferencesOperation;
    use NW\WebService\References\Operations\Notification\Views\MessageTemplateView;

    class ReturnOperationController extends ReferencesOperation
    {
        public function doOperation(): array {

            $request = new SendMessageRequest($this->getRequest('data'));
            $response = new SendNotificationResponse();

            if (!$request->resellerId) {
                $response->notificationClientBySms->message = 'Empty resellerId';

                return $response->toArray();
            }

            $notificationType = $request->notificationType;

            if (!$notificationType) {
                throw new \Exception('Empty notificationType', 400);
            }


            $reseller = Seller::getById($request->resellerId);
            if (!$reseller) {
                throw new \Exception('Seller not found!', 400);
            }


            $client = Contractor::getById($request->clientId);
            if ($client === null) throw new \Exception('client not found!', 400);


            //Если тип пользователя не клиент, или указанный продавец не закреплен за пользователем
            if ($client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $request->resellerId) {
                throw new \Exception('Bad client!', 400);
            }



            //$cr
            $cr = Employee::getById($request->creatorId);
            if ($cr === null) {
                throw new \Exception('Creator not found!', 400);
            }

            //$et
            $et = Employee::getById($request->expertId);
            if ($et === null) {
                throw new \Exception('Expert not found!', 400);
            }


            $emailFrom = getResellerEmailFrom($request->resellerId);
            if(empty($emailFrom) ){
                throw new \Exception('Email from not found!', 400);
            }


            $templateData = MessageTemplateView::getTemplate($request, $client, $cr, $et);

            $emailMessageTemplate= [];
            $emailMessageTemplate [] = [ // MessageTypes::EMAIL
                'emailFrom' => $emailFrom,
                'emailTo'   => null,
                'subject'   => __('complaintEmployeeEmailSubject', $templateData, $request->resellerId),
                'message'   => __('complaintEmployeeEmailBody', $templateData, $request->resellerId),
            ];


            $textReturnStatus =  NotificationStatusType::NEW_RETURN_STATUS;
            if ($notificationType == ENotificationDifferentType::TYPE_CHANGE){
                $textReturnStatus =  NotificationStatusType::CHANGE_RETURN_STATUS;
            }


            foreach (getEmailsByPermit( $request->resellerId, 'tsGoodsReturn') as $email){
                $emailMessageTemplate[0]['emailTo'] = $email;
                MessagesClient::sendMessage($emailMessageTemplate, $request->resellerId, $request->resellerId, $textReturnStatus, $request->differences->to);
                $response->notificationEmployeeByEmail = true;
            }


            if ($notificationType == ENotificationDifferentType::TYPE_CHANGE){
                $emailMessageTemplate[0]['emailTo'] = $client->email;
                MessagesClient::sendMessage($emailMessageTemplate, $request->resellerId, $client->id, $textReturnStatus, $request->differences->to);
                $response->notificationClientByEmail = true;


                if (!empty($client->mobile)) {
                    $error = null;
                    $res = NotificationManager::send($request->resellerId, $client->id, $textReturnStatus,  $request->differences->to, $templateData, $error);
                    if ($res) {
                        $response->notificationClientBySms->isSent = true;
                    }
                    if (!empty($error)) {
                        $response->notificationClientBySms->message = $error;
                    }
                }
            }


            return $response->toArray();
        }
    }
