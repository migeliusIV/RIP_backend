INSERT INTO public.gates (
    id_gate,
    title,
    description,
    status,
    image,
    full_info,
    the_axis
) VALUES
(
    1,
    'Identity Gate',
    'Не изменяет состояния кубита.',
    true,
    'http://127.0.0.1:9000/ibm-pictures/img/I-gate.png',
    'Ничего не делает с состоянием кубита. Оставляет его без изменений.',
    'non'
),
(
    2,
    'Pauli-X Gate (NOT gate)',
    'Инвертирует состояние кубита.',
    true,
    'http://127.0.0.1:9000/ibm-pictures/img/X-gate.png',
    'Аналог классического NOT-гейта. Переворачивает состояние кубита.',
    'non'
),
(
    3,
    'X-axis Rotation Gate',
    'Вращает кубит вокруг оси X на угол тэта.',
    true,
    'http://127.0.0.1:9000/ibm-pictures/img/X-rot-gate.png',
    'Эта операция вращает состояние кубита на сфере Блоха вокруг оси X. Значение угла поворота можно задать при компановке выражения (в деталях калькуляции).',
    'X'
),
(
    4,
    'Y-axis Rotation Gate',
    'Вращает кубит вокруг оси Y на угол тэта.',
    true,
    'http://127.0.0.1:9000/ibm-pictures/img/Y-rot-gate.png',
    'Эта операция вращает состояние кубита на сфере Блоха вокруг оси Y. Значение угла поворота можно задать при компановке выражения (в деталях калькуляции).',
    'Y'
),
(
    5,
    'Z-axis Rotation Gate',
    'Вращает кубит вокруг оси Z на угол тэта.',
    true,
    'http://127.0.0.1:9000/ibm-pictures/img/Z-rot-gate.png',
    'Эта операция вращает состояние кубита на сфере Блоха вокруг оси Z. Значение угла поворота можно задать при компановке выражения (в деталях калькуляции).',
    'Z'
),
(
    6,
    'H (Hadamard) Gate',
    'Создает равномерную суперпозицию из базисного состояния.',
    true,
    'http://127.0.0.1:9000/ibm-pictures/img/H-gate.png',
    'Операция поворачивает кубит на 90 градусов вокруг оси Y, затем на 180 градусов вокруг оси X. Это один из самых важных гейтов.',
    'non'
)

INSERT INTO public.degrees_to_gates (id_gate, id_task, degrees) VALUES
(1, 1),   -- Identity Gate (gates[0])
(2, 1),   -- Pauli-X Gate (gates[1]) 
(4, 1, 30);  -- Y-axis Rotation Gate (gates[3])

INSERT INTO public.quantum_tasks (
    id_task,
    task_status,
    creation_date,
    id_user,
    conclusion_date,
    id_moderator,
    task_description,
    res_koeff_0,
    res_koeff_1
) VALUES (
    1,
    'черновик',
    NOW(),
    1,
    NULL,
    NULL,
    'Компания АБВГД. Задача номер 1. Опровержение статистических гипотез.',
    0.2588,
    0.9659
);

-- Вставляем обычного пользователя
INSERT INTO public.users (login, password, is_admin) VALUES
('user1', 'password1', false),
('admin', 'adminpassword', true);