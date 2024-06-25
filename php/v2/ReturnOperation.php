<?php

namespace NW\WebService\References\Operations\Notification;

use DateTime;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;
    
    private static $default_result = [
        'notificationEmployeeByEmail' => false,
        'notificationClientByEmail'   => false,
        'notificationClientBySms'     => [
            'isSent'  => false,
            'message' => '',
        ],
    ];
    
    private $result;
    private $request_data;
    
    /**
     * @return array
     * @throws \Exception
     */
    public function doOperation(): array
    {
        // Prepare
        $this->request_data = $this->getRequest( 'data' );
        $this->result       = self::$default_result;
        
        // Validate request
        try{
            $this->validateRequestData( $this->request_data );
        }catch( \InvalidArgumentException $e ){
            throw new \Exception('Data validation error: ' . $e->getMessage(), 400 );
        }
        
        // Sanitize request
        $this->sanitizeRequestData( $this->request_data );
        
        // Generate actors
        $reseller = Seller::getById( (int)$this->request_data['resellerId'] );
        $client   = Client::getById( (int)$this->request_data['clientId'] );
        $creator  = Employee::getById( (int)$this->request_data['creatorId'] );
        $expert   = Employee::getById( (int)$this->request_data['expertId'] );
        
        // Email data
        $differences     = $this->getDifferences( $reseller );
        $template_data   = $this->createTemplateData( $client, $creator, $expert, $differences );
        $email_from      = getResellerEmailFrom( $reseller->getId() );
        $employee_emails = getEmailsByPermit( $reseller->getId(), 'tsGoodsReturn' );
        
        $this->sendEmailNotificationToEmployees(
            $employee_emails,
            $template_data,
            $email_from,
            $reseller
        );
        
        // Шлём клиентское уведомление, только если произошла смена статуса
        if( $this->request_data['notificationType'] === self::TYPE_CHANGE && ! empty( $this->request_data['differences']['to'] ) ){
            
            if( ! empty( $email_from ) && ! empty( $client->getEmail() ) ){
                $this->sendEmailNotificationToClients(
                    $client->getEmail(),
                    $template_data,
                    $email_from,
                    $reseller,
                    $this->request_data['differences']['to']
                );
            }
            
            if( ! empty( $client->getMobile() ) ){
                
                $error = false;
                $sms_notification_result = NotificationManager::send(
                    $reseller->getId(),
                    $client->getId(),
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    (int)$this->request_data['differences']['to'],
                    $template_data,
                    $error // @todo looks like $error passed by reference
                );
                
                if( $sms_notification_result ){
                    $this->result['notificationClientBySms']['isSent'] = true;
                }
                
                // @todo looks like $error passed by reference in NotificationManager::send
                if( ! empty( $error ) ){
                    $this->result['notificationClientBySms']['message'] = $error;
                }
                
            }
        }
        
        return $this->result;
    }
    
    /**
     * Validates request data
     *
     * @todo it could be done better with request schema and some framework builtin validator
     *       I decided not to complicate thing by wrote my own
     *
     * @param array $request_data
     *
     * @return void
     * @throws \Exception
     */
    private function validateRequestData( array $request_data )
    {
        // IDs
        if( $this->validateUID( $this->request_data['resellerId'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad resellerId', 400 );
        }
        if( $this->validateUID( $this->request_data['creatorId'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad creatorId', 400 );
        }
        if( $this->validateUID( $this->request_data['clientId'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad clientId', 400 );
        }
        if( $this->validateUID( $this->request_data['expertId'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad expertId', 400 );
        }
        if( $this->validateUID( $this->request_data['complaintId'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad complaintId', 400 );
        }
        if( $this->validateUID( $this->request_data['consumptionId'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad consumptionId', 400 );
        }
        
        // Validate "differences". Empty string or digits. Don't know what exactly needed
        if( preg_match( '^(|\d+)$', $this->request_data['differences']['to'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad differences TO', 400 );
        }
        if( preg_match( '^(|\d+)$', $this->request_data['differences']['from'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad differences TO', 400 );
        }
        
        // Numbers
        if( preg_match( '^(|\d+)$', $this->request_data['complaintNumber'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad complaintNumber', 400 );
        }
        if( preg_match( '^(|\d+)$', $this->request_data['consumptionNumber'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad consumptionNumber', 400 );
        }
        if( preg_match( '^(|\d+)$', $this->request_data['agreementNumber'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad agreementNumber', 400 );
        }
        
        // Date
        // @todo date format should be known for proper validation
        if( $this->validateDate( $this->request_data['date'] ?? '' ) ){
            throw new \InvalidArgumentException( 'Bad date', 400 );
        }
        
        // Status
        if( empty( (int)$this->request_data['notificationType'] ) ){
            throw new \InvalidArgumentException( 'Empty notificationType', 400 );
        }
    }
    
    /**
     * Validates UID and same IDs
     *
     * @param $id
     *
     * @return bool
     */
    private function validateUID( $id ): bool
    {
        return (bool) preg_match( '^[a-zA-Z0-9\-]{,128}$', $id );
    }
    
    /**
     * Validates date considering given format
     *
     * @param $date
     * @param $format
     *
     * @return bool
     */
    private function validateDate( $date, $format = 'Y-m-d H:i:s' )
    {
        $d = DateTime::createFromFormat( $format, $date );
        
        return $d && $d->format( $format ) === $date;
    }
    
    /**
     * Sanitizes request data
     *
     * @param array $request_data
     *
     * @return void
     */
    private function sanitizeRequestData( array &$request_data )
    {
        // @todo for proper sanitization, the allowed data range should be known
    }
    
    /**
     * Send email messages to employees
     *
     * @param array  $email_addresses
     * @param array  $template_data
     * @param string $email_from
     * @param Seller $reseller
     *
     * @return void
     */
    private function sendEmailNotificationToEmployees( $email_addresses, $template_data, $email_from, $reseller ): void
    {
        $email_addresses = (array) $email_addresses;
        
        foreach( $email_addresses as $email ){
            $this->sendEmailNotification( $email, $template_data, $email_from, $reseller );
        }
        
        $this->result['notificationEmployeeByEmail'] = true;
    }
    
    /**
     * Send email messages to employees
     *
     * @param array  $email_addresses
     * @param array  $template_data
     * @param string $email_from
     * @param Seller $reseller
     * @param mixed  $diff
     *
     * @return void
     */
    private function sendEmailNotificationToClients( $email_addresses, $template_data, $email_from, $reseller, $diff = null ): void
    {
        $email_addresses = (array) $email_addresses;
        
        foreach( $email_addresses as $email ){
            $this->sendEmailNotification( $email, $template_data, $email_from, $reseller, $diff );
        }
        
        $this->result['notificationClientByEmail'] = true;
    }
    
    /**
     * Send email messages
     *
     * @param string $email
     * @param array  $template_data
     * @param string $email_from
     * @param Seller $reseller
     * @param mixed  $diff @todo unknown parameter don't know what it's do need to see project context
     *
     * @return void
     */
    private function sendEmailNotification( $email, $template_data, $email_from, $reseller, $diff = null ): void
    {
        MessagesClient::sendMessage(
            [
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $email_from,
                    'emailTo'   => $email,
                    'subject'   => __( 'complaintEmployeeEmailSubject', $template_data, $reseller->getId() ),
                    'message'   => __( 'complaintEmployeeEmailBody',    $template_data, $reseller->getId() ),
                ],
            ],
            $reseller->getId(),
            NotificationEvents::CHANGE_RETURN_STATUS,
            $diff
        );
    }
    
    /**
     * Compile template data for notification
     *
     * @param Client   $client
     * @param Employee $creator
     * @param Employee $expert
     * @param array    $differences
     *
     * @return array
     * @throws \Exception
     */
    private function createTemplateData( $client, $creator, $expert, $differences ): array
    {
        $templateData = [
            'COMPLAINT_ID'       => (int)$this->request_data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$this->request_data['complaintNumber'],
            
            'CREATOR_ID'         => $creator->getId(),
            'CREATOR_NAME'       => $creator->getFullName(),
            
            'EXPERT_ID'          => $expert->getId(),
            'EXPERT_NAME'        => $expert->getFullName(),
            
            'CLIENT_ID'          => $client->getId(),
            'CLIENT_NAME'        => $client->getFullName(),
            
            'CONSUMPTION_ID'     => (int)$this->request_data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$this->request_data['consumptionNumber'],
            
            'AGREEMENT_NUMBER'   => (string)$this->request_data['agreementNumber'],
            'DATE'               => (string)$this->request_data['date'],
            'DIFFERENCES'        => $differences,
        ];
        
        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach( $templateData as $key => $tempData ){
            if( empty( $tempData ) ){
                throw new \Exception( "Template Data ({$key}) is empty!", 500 );
            }
        }
        
        return $templateData;
    }
    
    /**
     * Get differences ( no matter what it means )
     *
     * @param Seller $reseller
     *
     * @return string
     * @throws \Exception
     */
    private function getDifferences( $reseller ): string
    {
        if( $this->request_data['notificationType'] === self::TYPE_NEW ){
            $differences = __( 'NewPositionAdded', null, $reseller->getId() );
            
        }elseif( $this->request_data['notificationType'] === self::TYPE_CHANGE && ! empty( $this->request_data['differences'] ) ){
            $differences = __(
                'PositionStatusHasChanged',
                [
                    'FROM' => Status::getNameById( (int)$this->request_data['differences']['from'] ),
                    'TO'   => Status::getNameById( (int)$this->request_data['differences']['to'] ),
                ],
                $reseller->getId()
            );
            
        }else{
            throw new \Exception( 'No difference occur, shutdown notifications', 400 );
        }
        
        return $differences;
    }

}
