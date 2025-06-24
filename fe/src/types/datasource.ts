import { ApiDataSourceResponse } from '@/services';

export type DataSource = NonNullable<ApiDataSourceResponse['list']>[0];
