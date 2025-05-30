CREATE TABLE exam_category (
                               id INT PRIMARY KEY AUTO_INCREMENT,
                               name VARCHAR(50) NOT NULL COMMENT '考试类别名称，如"系统架构设计师"',
                               level VARCHAR(20) COMMENT '初级/中级/高级',
                               description TEXT
);

CREATE TABLE exam_question (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
                          category_id INT NOT NULL COMMENT '考试类别ID',
                          question_type TINYINT NOT NULL COMMENT '1-单选 2-多选 3-判断 4-案例 5-论文',
                          content TEXT NOT NULL COMMENT '题目内容',
                          year SMALLINT COMMENT '年份',
                          created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          INDEX idx_category (category_id) COMMENT '类别索引',
                          INDEX idx_year (year) COMMENT '年份索引'
);

CREATE TABLE exam_question_option (
                                 id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
                                 question_id BIGINT NOT NULL COMMENT '题目ID，逻辑关联question.id',
                                 option_key CHAR(1) NOT NULL COMMENT '选项键，如A/B/C/D',
                                 content TEXT NOT NULL COMMENT '选项内容',
                                 is_correct BOOLEAN DEFAULT FALSE COMMENT '是否正确答案',
                                 INDEX idx_question (question_id) COMMENT '题目ID索引'
);

CREATE TABLE exam_question_analysis (
                                   id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
                                   question_id BIGINT NOT NULL UNIQUE COMMENT '题目ID，逻辑关联question.id',
                                   analysis TEXT COMMENT '题目解析',
                                   knowledge_points TEXT COMMENT '涉及知识点',
                                   INDEX idx_question (question_id) COMMENT '题目ID索引'
);

CREATE TABLE exam_knowledge_point (
                                 id INT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
                                 category_id INT NOT NULL COMMENT '考试类别ID，逻辑关联exam_category.id',
                                 name VARCHAR(100) NOT NULL COMMENT '知识点名称',
                                 INDEX idx_category (category_id) COMMENT '类别索引'
);

CREATE TABLE exam_question_knowledge (
                                    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
                                    question_id BIGINT NOT NULL COMMENT '题目ID，逻辑关联question.id',
                                    knowledge_id INT NOT NULL COMMENT '知识点ID，逻辑关联knowledge_point.id',
                                    INDEX idx_question (question_id) COMMENT '题目ID索引',
                                    INDEX idx_knowledge (knowledge_id) COMMENT '知识点ID索引',
                                    UNIQUE KEY uk_question_knowledge (question_id, knowledge_id) COMMENT '联合唯一约束'
);

CREATE TABLE exam_question_favorite (
                                    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
                                    user_id BIGINT NOT NULL COMMENT '用户ID，逻辑关联user.id',
                                    question_id BIGINT NOT NULL COMMENT '题目ID，逻辑关联question.id',
                                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    INDEX idx_user (user_id) COMMENT '用户ID索引',
                                    INDEX idx_question (question_id) COMMENT '题目ID索引'
)

CREATE TABLE exam_user (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
                           nickname VARCHAR(50) NOT NULL COMMENT '用户昵称',
                           avatar VARCHAR(200) COMMENT '用户头像',
                           account VARCHAR(50) NOT NULL COMMENT '用户账号',
                           password VARCHAR(100) NOT NULL COMMENT '用户密码',
                           created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
)