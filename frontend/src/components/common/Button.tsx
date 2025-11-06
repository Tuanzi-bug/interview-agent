import { Button as AntButton, ButtonProps as AntButtonProps } from 'antd';

export interface ButtonProps extends AntButtonProps {
  variant?: 'primary' | 'default' | 'dashed' | 'text' | 'link';
}

const Button: React.FC<ButtonProps> = ({ variant = 'default', ...props }) => {
  return <AntButton type={variant === 'primary' ? 'primary' : 'default'} {...props} />;
};

export default Button;