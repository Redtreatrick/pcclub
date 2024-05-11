Хочу ли я обрабатывать события асинхронно?

3
09:00 19:00
10
08:48 1 client1
09:41 1 client1
09:48 1 client2
09:52 3 client1
09:54 2 client1 1 // client1 сел за стол 1 в 09:54
10:25 2 client2 2
10:58 1 client3
10:59 2 client3 3
11:30 1 client4
11:35 2 client4 2
11:45 3 client4 // client4 ждёт в очереди 
12:33 4 client1 // client1 вышел из-за стола 1 в 12:33, вместо него сел client4
12:43 4 client2
15:52 4 client4 // client4 вышел из-за стола 1 в 15:52