cub kafka-ready -b kafka:9092 1 20

sleep 5

topics=("friends_operation" "session_SET" "session_UPDATE" "push_friends" "push_subscribers" "update_subscriber" "update_friend", "message_operation", "read_notification")

for topic in "${topics[@]}"; do
    kafka-topics --create --topic "$topic" --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1
    if [ $? -eq 0 ]; then
        echo "Топик '$topic' успешно создан."
    else
        echo "Ошибка при создании топика '$topic'."
    fi
done