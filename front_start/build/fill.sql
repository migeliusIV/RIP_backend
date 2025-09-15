INSERT INTO public.gates (id_gate, title, description, status, image, full_info, is_editable, the_axis) VALUES
(1, 'Identity Gate', 'Не изменяет состояния кубита.', true, 'http://127.0.0.1:9000/ibm-pictures/img/I-gate.png', 'Ничего не делает с состоянием кубита. Оставляет его без изменений.', false, ''),
(2, 'Pauli-X Gate (NOT gate)', 'Инвертирует состояние кубита.', true, 'http://127.0.0.1:9000/ibm-pictures/img/X-gate.png', 'Аналог классического NOT-гейта. Переворачивает состояние кубита.', false, ''),
(3, 'X-axis Rotation Gate', 'Вращает кубит вокруг оси X на угол тэта.', true, 'http://127.0.0.1:9000/ibm-pictures/img/X-rot-gate.png', 'Эта операция вращает состояние кубита на сфере Блоха вокруг оси X. Значение угла поворота можно задать при компановке выражения (в деталях калькуляции).', true, 'X'),
(4, 'Y-axis Rotation Gate', 'Вращает кубит вокруг оси Y на угол тэта.', true, 'http://127.0.0.1:9000/ibm-pictures/img/Y-rot-gate.png', 'Эта операция вращает состояние кубита на сфере Блоха вокруг оси Y. Значение угла поворота можно задать при компановке выражения (в деталях калькуляции).', true, 'Y'),
(5, 'Z-axis Rotation Gate', 'Вращает кубит вокруг оси Z на угол тэта.', true, 'http://127.0.0.1:9000/ibm-pictures/img/Z-rot-gate.png', 'Эта операция вращает состояние кубита на сфере Блоха вокруг оси Z. Значение угла поворота можно задать при компановке выражения (в деталях калькуляции).', true, 'Z'),
(6, 'H (Hadamard) Gate', 'Создает равномерную суперпозицию из базисного состояния.', true, 'http://127.0.0.1:9000/ibm-pictures/img/H-gate.png', 'Операция поворачивает кубит на 90 градусов вокруг оси Y, затем на 180 градусов вокруг оси X. Это один из самых важных гейтов.', false, '');

INSERT INTO public.degrees_to_gates (id_gate, id_task, degrees) VALUES
(1, 1, 0),   -- Identity Gate (gates[0])
(2, 1, 0),   -- Pauli-X Gate (gates[1]) 
(4, 1, 30);  -- Y-axis Rotation Gate (gates[3])

INSERT INTO public.tasks (
    id_task, 
    tesk_status, 
    creation_date, 
    id_user, 
    conclusion_date, 
    id_moderator, 
    task_description, 
    result
) VALUES (
    1, 
    'черновик', 
    NOW(), 
    1,  -- предполагая, что пользователь с id_user=1 существует
    NULL, 
    NULL, 
    'Компания АБВГД. Задача номер 1. Опровержение статистических гипотез.', 
    '0.2588|0⟩ + 0.9659|1⟩'
);

-- Вставляем обычного пользователя
INSERT INTO public.users (login, password, is_admin) VALUES
('user1', 'password1', false),
('admin', 'adminpassword', true);