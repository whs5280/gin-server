-- 插入考试类别
INSERT INTO exam_category (id, name, level, description) VALUES (1, '软件设计师', '中级', '面向软件开发的中级资格考试,https://www.ruankao.org.cn/');

-- 插入题目
INSERT INTO exam_question (id, category_id, question_type, content, year) VALUES
      (1, 1, 1, '在面向对象设计中，以下哪个原则建议"对扩展开放，对修改关闭"？', 2023),
      (2, 1, 1, '以下哪种排序算法在最坏情况下的时间复杂度是O(n^2)？', 2023),
      (3, 1, 2, '关于数据库事务的ACID特性，以下描述正确的是？', 2023),
      (4, 1, 1, '当需要为一个对象动态添加额外职责时，应该使用哪种设计模式？', 2023),
      (5, 1, 1, '在UML类图中，下列哪种关系表示"整体-部分"关系？', 2023),
      (6, 1, 1, '以下哪种测试主要关注系统是否符合用户需求？', 2023),
      (7, 1, 1, '下列哪种进程调度算法可能导致"饥饿"现象？', 2023),
      (8, 1, 1, 'CPU中的哪个部件负责执行算术和逻辑运算？', 2023);

-- 插入选项
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (1, 'A', '单一职责原则', false),
                                                                               (1, 'B', '开闭原则', true),
                                                                               (1, 'C', '里氏替换原则', false),
                                                                               (1, 'D', '接口隔离原则', false);
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (2, 'A', '归并排序', false),
                                                                               (2, 'B', '快速排序', true),
                                                                               (2, 'C', '堆排序', false),
                                                                               (2, 'D', '希尔排序', false);
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (3, 'A', '原子性指事务是不可分割的工作单位', true),
                                                                               (3, 'B', '一致性指事务执行前后数据必须完全一致', false),
                                                                               (3, 'C', '隔离性指多个事务并发执行时互不干扰', true),
                                                                               (3, 'D', '持久性指事务提交后对系统的影响是永久的', true);
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (4, 'A', '策略模式', false),
                                                                               (4, 'B', '装饰器模式', true),
                                                                               (4, 'C', '工厂模式', false),
                                                                               (4, 'D', '观察者模式', false);
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (5, 'A', '关联关系', false),
                                                                               (5, 'B', '聚合关系', true),
                                                                               (5, 'C', '泛化关系', false),
                                                                               (5, 'D', '实现关系', false);
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (6, 'A', '单元测试', false),
                                                                               (6, 'B', '集成测试', false),
                                                                               (6, 'C', '系统测试', false),
                                                                               (6, 'D', '验收测试', true);
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (7, 'A', '时间片轮转', false),
                                                                               (7, 'B', '短作业优先', true),
                                                                               (7, 'C', '先来先服务', false),
                                                                               (7, 'D', '多级反馈队列', false);
INSERT INTO exam_question_option (question_id, option_key, content, is_correct) VALUES
                                                                               (8, 'A', '控制器', false),
                                                                               (8, 'B', '运算器', true),
                                                                               (8, 'C', '寄存器', false),
                                                                               (8, 'D', '缓存', false);

-- 插入知识点（软件设计师）
INSERT INTO knowledge_point (id, category_id, name, parent_id, level) VALUES
                                                                          (1, 1, '面向对象设计', NULL, 1),
                                                                          (2, 1, '算法与数据结构', NULL, 1),
                                                                          (3, 1, '数据库设计', NULL, 1),
                                                                          (4, 1, '设计模式', 6, 2),
                                                                          (5, 1, 'UML建模', 6, 2),
                                                                          (6, 1, '软件测试', NULL, 1),
                                                                          (7, 1, '操作系统', NULL, 1),
                                                                          (8, 1, '计算机组成原理', NULL, 1),
                                                                          (9, 1, '软件工程', NULL, 1),
                                                                          (10, 1, '编程语言基础', NULL, 1);
