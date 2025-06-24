import { Tooltip } from 'antd';
import dayjs from 'dayjs';
import { isNaN } from 'lodash';

export const formatDuration = (text: string | number) => {
  const num = +text;
  let str = '-';
  let tooltip = '';
  const duration = dayjs.duration(num, 's');
  const _f = (str: string) => duration.format(str);
  if (isNaN(num)) {
    str = '-';
  } else if (duration.asYears() >= 1) {
    str = _f('y[Y]');
    tooltip = _f('y [years] M [months]');
  } else if (duration.asMonths() >= 1) {
    str = _f('M[M]');
    tooltip = _f('M [months] d [days]');
  } else if (duration.asDays() >= 1) {
    str = _f('D[D]');
    tooltip = _f('D [days] H [hours]');
  } else if (duration.asHours() >= 1) {
    str = _f('H[h]');
    tooltip = _f('H [hours] m [minutes]');
  } else if (duration.asMinutes() >= 1) {
    str = _f('m[m]');
    tooltip = _f('m [minutes] s [seconds]');
  } else {
    str = _f('s[s]');
    tooltip = _f('s [seconds]');
  }
  return (
    <Tooltip title={tooltip}>
      <span>{str}</span>
    </Tooltip>
  );
};
