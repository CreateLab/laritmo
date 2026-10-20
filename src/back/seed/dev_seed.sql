-- Тестовые данные для локальной разработки. НЕ применять на проде.
-- Повторный запуск пересоздаёт курсы и всё, что к ним привязано.
SET NAMES utf8mb4;
DELETE FROM courses;
ALTER TABLE courses AUTO_INCREMENT = 1;
ALTER TABLE lectures AUTO_INCREMENT = 1;
ALTER TABLE labs AUTO_INCREMENT = 1;
ALTER TABLE grade_sheets AUTO_INCREMENT = 1;
ALTER TABLE exam_questions AUTO_INCREMENT = 1;

INSERT INTO courses (name, semester, description) VALUES
('Программирование (Си)', '1 семестр', 'Основы языка Си: типы данных, указатели, работа с памятью и файлами.'),
('Основы ASP.NET Core', '5 семестр', 'Разработка Web API: маршрутизация, middleware, работа с базой данных.'),
('Базы данных', '4 семестр', 'Реляционная модель, SQL, проектирование схем и транзакции.'),
('Алгоритмы и структуры данных', '3 семестр', 'Сложность, сортировки, деревья, графы и динамическое программирование.');

INSERT INTO lectures (course_id, week, title, content, github_url) VALUES
(1, 1, 'Введение в язык Си', 'Язык Си создан в 1970-х годах Деннисом Ритчи.\n\n## Первая программа\n\n```c\n#include <stdio.h>\n\nint main(void) {\n    printf("Hello, world!\\n");\n    return 0;\n}\n```\n\nКомпиляция: `gcc main.c -o main`.\n\n## Что дальше\n\n- Типы данных\n- Операторы\n- Функции\n\nПодробнее в [документации](https://en.cppreference.com/w/c).', 'https://github.com/example/c-course/tree/main/week1'),
(1, 2, 'Указатели и массивы', 'Указатель хранит адрес переменной.\n\n```c\nint x = 42;\nint *p = &x;\nprintf("%d\n", *p);\n```\n\n## Массивы\n\nИмя массива — указатель на его первый элемент.\n\n| Выражение | Смысл |\n| --- | --- |\n| `a[i]` | i-й элемент |\n| `*(a + i)` | то же самое |', NULL),
(1, 3, 'Динамическая память', 'Функции `malloc` и `free` выделяют и освобождают память в куче.\n\n```c\nint *a = malloc(10 * sizeof(int));\nif (a == NULL) return 1;\nfree(a);\n```\n\n> Каждому `malloc` соответствует ровно один `free`.', NULL),
(2, 1, 'Структура проекта ASP.NET Core', 'Проект содержит `Program.cs`, `appsettings.json` и папку `Controllers`.\n\n```csharp\nvar builder = WebApplication.CreateBuilder(args);\nvar app = builder.Build();\napp.MapGet("/", () => "ok");\napp.Run();\n```', 'https://github.com/example/aspnet/tree/main/week1'),
(2, 2, 'Маршрутизация и контроллеры', 'REST — это архитектурный стиль.\n\n## Атрибуты\n\n```csharp\n[ApiController]\n[Route("api/[controller]")]\npublic class CoursesController : ControllerBase\n{\n    [HttpGet]\n    public IActionResult Get() => Ok();\n}\n```', NULL),
(3, 1, 'Реляционная модель', 'Данные хранятся в таблицах. Строка — кортеж, столбец — атрибут.\n\n```sql\nSELECT name, semester FROM courses WHERE semester = ''1 семестр'';\n```', NULL);

INSERT INTO labs (course_id, number, title, description, deadline, max_score, github_url) VALUES
(1, 1, 'Калькулятор в консоли', 'Напишите консольный калькулятор для четырёх арифметических операций.\n\n## Требования\n\n- Ввод через `scanf`\n- Обработка деления на ноль', DATE_ADD(NOW(), INTERVAL 14 DAY), 100, 'https://github.com/example/c-course/tree/main/lab1'),
(1, 2, 'Односвязный список', 'Реализуйте односвязный список с операциями добавления, удаления и поиска.', DATE_ADD(NOW(), INTERVAL 30 DAY), 100, NULL),
(2, 1, 'REST API для курсов', 'Создайте Web API с операциями CRUD для сущности «Курс».', DATE_SUB(NOW(), INTERVAL 7 DAY), 50, NULL);

INSERT INTO grade_sheets (course_id, sheet_url, description) VALUES
(1, 'https://docs.google.com/spreadsheets/d/example-sheet-id/edit', 'Журнал группы 101');

INSERT INTO exam_questions (course_id, number, section, question) VALUES
(1, 1, 'Основы', 'Чем отличается объявление переменной от её определения?'),
(1, 2, 'Основы', 'Перечислите базовые типы данных языка Си и их размеры.'),
(1, 3, 'Указатели', 'Что такое указатель и как работает операция разыменования?'),
(1, 4, 'Указатели', 'Как связаны массивы и указатели в языке Си?'),
(1, 5, 'Память', 'Чем стек отличается от кучи? Что такое утечка памяти?'),
(2, 1, 'Основы', 'Что такое middleware в ASP.NET Core?'),
(2, 2, 'Основы', 'Как работает внедрение зависимостей в ASP.NET Core?');
