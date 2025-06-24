import { Button, message } from 'antd';
import { useCopyToClipboard } from 'react-use';
import { CopyOutlined } from '@ant-design/icons';

const EMAIL_REGEX = /^([a-zA-Z0-9_.+-]+)@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$/;
export const RenderEmail = ({ value }: { value: unknown }) => {
  const [_, copyToClipboard] = useCopyToClipboard();
  const onCopy = (value: string) => {
    copyToClipboard(value);
    message.success(`Copied text: ${value}`);
  };
  if (!value || typeof value !== 'string') return 'N/A';
  if (!EMAIL_REGEX.test(value)) return value;
  const prefix = value.split('@')[0];
  return (
    <span>
      <Button
        icon={<CopyOutlined />}
        size="small"
        type="text"
        onClick={() => onCopy(value)}
      ></Button>
      {prefix}
    </span>
  );
};
