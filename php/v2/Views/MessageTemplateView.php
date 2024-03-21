<?php

    namespace NW\WebService\References\Operations\Notification\Views;



    use NW\WebService\References\Operations\Notification\BusinessLayer\Enums\ENotificationDifferentType;
    use NW\WebService\References\Operations\Notification\BusinessLayer\SendMessageRequest;
    use NW\WebService\References\Operations\Notification\Contractor;

    class MessageTemplateView{


        public static function getTemplate(SendMessageRequest $request, Contractor $client, Contractor $cr, Contractor $et){

            $notificationType = $request->notificationType;

            $differences = '';
            if ($notificationType == ENotificationDifferentType::TYPE_NEW) {
                $differences = __('NewPositionAdded', null, $request->resellerId);
            } elseif ($notificationType == ENotificationDifferentType::TYPE_CHANGE && !empty($data['differences'])) {
                $differences = __('PositionStatusHasChanged', [
                    'FROM' => $request->differences->getFromName(),
                    'TO'   =>$request->differences->getToName(),
                ], $request->resellerId);
            }

            $templateData = [
                'COMPLAINT_ID'       => $request->complaintId,
                'COMPLAINT_NUMBER'   => $request->complaintNumber,
                'CREATOR_ID'         => $request->creatorId,
                'CREATOR_NAME'       => $cr->getFullName(),
                'EXPERT_ID'          => $request->expertId,
                'EXPERT_NAME'        => $et->getFullName(),
                'CLIENT_ID'          => $request->clientId,
                'CLIENT_NAME'        => $client->getFullName() ?? $client->name,  //вместо $cFullName
                'CONSUMPTION_ID'     => $request->consumptionId,
                'CONSUMPTION_NUMBER' => $request->consumptionNumber,
                'AGREEMENT_NUMBER'   => $request->agreementNumber,
                'DATE'               => $request->date,
                'DIFFERENCES'        => $differences,
            ];

            // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
            foreach ($templateData as $key => $tempData) {
                if (empty($tempData)) {
                    throw new \Exception("Template Data ({$key}) is empty!", 500);
                }
            }


            return $templateData;
        }

    }
