import { Card, Button, Space, Tag, Progress, Typography, Row, Col } from 'antd';
import { PlayCircleOutlined, ClockCircleOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

const { Title, Text, Paragraph } = Typography;

interface SpecialInterviewCardProps {
  title: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  estimatedTime: number;
  completedQuestions: number;
  totalQuestions: number;
  icon: React.ReactNode;
  color: string;
  onStart: () => void;
}

const SpecialInterviewCard: React.FC<SpecialInterviewCardProps> = ({
  title,
  description,
  difficulty,
  estimatedTime,
  completedQuestions,
  totalQuestions,
  icon,
  color,
  onStart
}) => {
  const progress = (completedQuestions / totalQuestions) * 100;
  
  const difficultyColor = {
    beginner: 'green',
    intermediate: 'orange',
    advanced: 'red'
  };

  const difficultyText = {
    beginner: '初级',
    intermediate: '中级',
    advanced: '高级'
  };

  return (
    <Card
      hoverable
      className="special-interview-card"
      style={{ 
        borderTop: `4px solid ${color}`,
        height: '100%',
        display: 'flex',
        flexDirection: 'column'
      }}
      bodyStyle={{ 
        flex: 1, 
        display: 'flex', 
        flexDirection: 'column',
        justifyContent: 'space-between'
      }}
    >
      <div>
        <div className="flex items-center mb-4">
          <div className="text-2xl mr-3" style={{ color }}>
            {icon}
          </div>
          <Title level={4} className="!mb-0">{title}</Title>
        </div>
        
        <Paragraph className="text-gray-600 mb-4" ellipsis={{ rows: 2 }}>
          {description}
        </Paragraph>
        
        <Space className="mb-4" wrap>
          <Tag color={difficultyColor[difficulty]}>
            {difficultyText[difficulty]}
          </Tag>
          <Tag icon={<ClockCircleOutlined />}>
            预计{estimatedTime}分钟
          </Tag>
        </Space>
        
        {completedQuestions > 0 && (
          <div className="mb-4">
            <div className="flex justify-between items-center mb-2">
              <Text className="text-sm text-gray-600">进度</Text>
              <Text className="text-sm text-gray-600">
                {completedQuestions}/{totalQuestions}
              </Text>
            </div>
            <Progress 
              percent={progress} 
              size="small" 
              strokeColor={color}
              format={() => `${Math.round(progress)}%`}
            />
          </div>
        )}
      </div>
      
      <Button
        type="primary"
        size="large"
        icon={<PlayCircleOutlined />}
        onClick={onStart}
        className="w-full"
        style={{ backgroundColor: color, borderColor: color }}
      >
        {completedQuestions > 0 ? '继续面试' : '开始面试'}
      </Button>
    </Card>
  );
};

export default function GoSpecialInterviewPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const handleStartInterview = async (category: string, difficulty: string) => {
    setLoading(true);
    try {
      // 调用后端API创建面试会话
      const response = await fetch('/api/interview/special/start', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          domain: 'go',
          category,
          difficulty,
          sessionType: 'special'
        })
      });

      if (!response.ok) {
        throw new Error('Failed to start interview');
      }

      const data = await response.json();
      
      // 跳转到面试页面
      router.push(`/interview/special/session/${data.sessionId}`);
    } catch (error) {
      console.error('Start interview error:', error);
      // 这里可以添加错误提示
    } finally {
      setLoading(false);
    }
  };

  const interviewCategories = [
    {
      key: 'fundamentals',
      title: 'Go基础语法',
      description: '掌握Go语言的基本语法、数据类型、控制结构等核心概念，为后续深入学习打下坚实基础',
      difficulty: 'beginner' as const,
      estimatedTime: 15,
      completedQuestions: 3,
      totalQuestions: 20,
      icon: '📝',
      color: '#52c41a'
    },
    {
      key: 'concurrency',
      title: '并发编程',
      description: '深入理解Goroutine、Channel、并发模式等Go语言的核心特性，掌握高并发程序设计',
      difficulty: 'intermediate' as const,
      estimatedTime: 25,
      completedQuestions: 0,
      totalQuestions: 15,
      icon: '⚡',
      color: '#fa8c16'
    },
    {
      key: 'memory',
      title: '内存管理',
      description: '了解Go语言的内存分配机制、垃圾回收原理，掌握内存优化技巧和性能调优方法',
      difficulty: 'advanced' as const,
      estimatedTime: 30,
      completedQuestions: 0,
      totalQuestions: 12,
      icon: '💾',
      color: '#f5222d'
    },
    {
      key: 'stdlib',
      title: '标准库',
      description: '熟悉Go标准库的核心包，如fmt、io、net、http等，提高开发效率和代码质量',
      difficulty: 'intermediate' as const,
      estimatedTime: 20,
      completedQuestions: 2,
      totalQuestions: 18,
      icon: '📚',
      color: '#1890ff'
    },
    {
      key: 'performance',
      title: '性能优化',
      description: '学习Go程序的性能分析工具和方法，掌握代码优化技巧，构建高性能应用',
      difficulty: 'advanced' as const,
      estimatedTime: 35,
      completedQuestions: 0,
      totalQuestions: 10,
      icon: '🚀',
      color: '#722ed1'
    },
    {
      key: 'testing',
      title: '测试开发',
      description: '掌握Go语言的测试框架和最佳实践，编写高质量的单元测试和集成测试',
      difficulty: 'intermediate' as const,
      estimatedTime: 18,
      completedQuestions: 0,
      totalQuestions: 16,
      icon: '🧪',
      color: '#13c2c2'
    }
  ];

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* 页面标题 */}
        <div className="text-center mb-12">
          <Title level={1} className="!mb-4">
            Go语言专项面试
          </Title>
          <Paragraph className="text-lg text-gray-600 max-w-3xl mx-auto">
            针对Go语言的专项技术面试，涵盖基础语法、并发编程、内存管理等核心技术领域
          </Paragraph>
        </div>

        {/* 统计信息 */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-8">
          <Row gutter={24}>
            <Col xs={24} sm={8}>
              <div className="text-center">
                <div className="text-3xl font-bold text-blue-600">6</div>
                <Text className="text-gray-600">专项领域</Text>
              </div>
            </Col>
            <Col xs={24} sm={8}>
              <div className="text-center">
                <div className="text-3xl font-bold text-green-600">91</div>
                <Text className="text-gray-600">面试题目</Text>
              </div>
            </Col>
            <Col xs={24} sm={8}>
              <div className="text-center">
                <div className="text-3xl font-bold text-orange-600">5</div>
                <Text className="text-gray-600">已完成</Text>
              </div>
            </Col>
          </Row>
        </div>

        {/* 专项面试卡片 */}
        <Row gutter={[24, 24]}>
          {interviewCategories.map((category) => (
            <Col key={category.key} xs={24} sm={12} lg={8}>
              <SpecialInterviewCard
                title={category.title}
                description={category.description}
                difficulty={category.difficulty}
                estimatedTime={category.estimatedTime}
                completedQuestions={category.completedQuestions}
                totalQuestions={category.totalQuestions}
                icon={category.icon}
                color={category.color}
                onStart={() => handleStartInterview(category.key, category.difficulty)}
              />
            </Col>
          ))}
        </Row>

        {/* 学习路径建议 */}
        <div className="mt-12 bg-white rounded-lg shadow-sm p-6">
          <Title level={3} className="!mb-6">学习路径建议</Title>
          <Row gutter={24}>
            <Col xs={24} md={8}>
              <div className="text-center p-4">
                <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-green-600 text-xl">1</span>
                </div>
                <Title level={4} className="!mb-2">基础阶段</Title>
                <Text className="text-gray-600">
                  从Go基础语法开始，掌握语言核心概念和基本用法
                </Text>
              </div>
            </Col>
            <Col xs={24} md={8}>
              <div className="text-center p-4">
                <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-orange-600 text-xl">2</span>
                </div>
                <Title level={4} className="!mb-2">进阶阶段</Title>
                <Text className="text-gray-600">
                  深入学习并发编程和标准库，提升实际开发能力
                </Text>
              </div>
            </Col>
            <Col xs={24} md={8}>
              <div className="text-center p-4">
                <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-red-600 text-xl">3</span>
                </div>
                <Title level={4} className="!mb-2">高级阶段</Title>
                <Text className="text-gray-600">
                  掌握性能优化和底层原理，成为Go专家
                </Text>
              </div>
            </Col>
          </Row>
        </div>
      </div>
    </div>
  );
}