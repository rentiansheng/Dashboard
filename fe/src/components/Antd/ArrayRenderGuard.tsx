import { ConfigContextPropsType, ProRenderFieldPropsType } from '@ant-design/pro-components';
import 'dayjs/locale/zh-cn';
import { EllipsisRender } from './EllipsisRender';
import styled from 'styled-components';

export const ArrayRenderGuard =
  (
    valueType: string,
    valueMap: ConfigContextPropsType['valueTypeMap'],
  ): ProRenderFieldPropsType['render'] =>
  (...args) => {
    const value = args[0];
    const renderText = valueMap?.[valueType]?.render;
    const render = (text: unknown) => {
      return text !== null && text !== undefined ? (
        <span>{renderText?.apply(null, [text, args[1], <></>]) || <>{text.toString()}</>}</span>
      ) : (
        <>N/A</>
      );
    };
    return (
      <EllipsisRender ellipsis={args[1]?.fieldProps?.ellipsis}>
        {Array.isArray(value) ? (
          <>
            {value.map((text) => (
              <>
                <BulletedSpan noWrap={args[1]?.fieldProps?.noWrap}>{render(text)}</BulletedSpan>
                <br />
              </>
            ))}
          </>
        ) : (
          render(value)
        )}
      </EllipsisRender>
    );
  };

const BulletedSpan = styled.span<{ noWrap?: boolean }>`
  &:before {
    content: '•';
    margin-right: 4px;
    font-size: 20px;
    line-height: 1;
  }
`;
