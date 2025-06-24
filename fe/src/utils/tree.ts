import { cloneDeep, omit } from 'lodash';
import { Required } from 'utility-types';

export type TraverseTreeOptions = {
  childrenKey?: string;
  by?: 'dfs' | 'bfs';
};

export type MapTreeOptions<T, R> = {
  parent?: R;
  children?: T[];
};

const resolveOptions = (options?: TraverseTreeOptions): Required<TraverseTreeOptions> => {
  return Object.assign(
    {
      childrenKey: 'children',
      by: 'dfs',
    },
    options,
  );
};

export function traverseTree<T>(
  root: T,
  cb: (node: T, parent: T | undefined) => void,
  options?: TraverseTreeOptions,
) {
  const { childrenKey, by } = resolveOptions(options);
  const isDFS = by === 'dfs';
  let pending: [T, T | undefined][] = [[root, void 0]];
  while (pending.length) {
    const [node, parent] = isDFS ? pending.pop()! : pending.shift()!;
    cb(node, parent);
    if ((node as any)[childrenKey] && (node as any)[childrenKey].length) {
      const list = (node as any)[childrenKey].map((child: T) => [child, node]);
      pending = pending.concat(isDFS ? list.reverse() : list);
    }
  }
}

export function mapTree<T, R>(
  root: T,
  cb: (node: T, options: MapTreeOptions<T, R>) => R,
  options?: Omit<TraverseTreeOptions, 'by'>,
): R | undefined {
  const { childrenKey } = resolveOptions(options);
  if (!root) return undefined;
  let newRoot: R | undefined;
  let pending: [T, R | undefined][] = [[root, newRoot]];
  while (pending.length) {
    const [node, newParent] = pending.pop()!;
    const newNode: any = cb((omit as any)(node, [childrenKey]), {
      parent: newParent,
      children: (node as any)[childrenKey],
    });
    const hasChildren = !!(node as any)[childrenKey]?.length;
    if (hasChildren) {
      newNode[childrenKey] = newNode[childrenKey] || [];
    }
    if (!newParent) {
      newRoot = newNode;
    } else {
      (newParent as any)[childrenKey].push(newNode);
    }
    if (hasChildren) {
      const list = (node as any)[childrenKey].map((child: T) => {
        return [child, newNode];
      });
      pending = pending.concat(list.reverse());
    }
  }
  return newRoot as R | undefined;
}

export function flattenTree<T>(root: T, childrenKey = 'children') {
  const nodes: T[] = [];
  traverseTree(
    root,
    (item) => {
      nodes.push(item);
    },
    { childrenKey },
  );
  return nodes;
}

export function traverseTreeList<T>(
  treeList: T[],
  cb: (node: T, parent: T | undefined) => void,
  options?: TraverseTreeOptions,
) {
  treeList.forEach((item) => traverseTree(item, cb, options));
}

export function flattenTreeList<T>(treeList: T[], childrenKey = 'children'): T[] {
  return treeList.reduce((acc, node) => {
    return acc.concat(flattenTree(node, childrenKey));
  }, [] as T[]);
}

export interface FilterTreeOptions<T> extends TraverseTreeOptions {
  filterChildren?: boolean;
  filterStrictly?: boolean;
  onlyLeaf?: boolean;
  isLeaf?: (node: T) => boolean;
}

function filterTreeListStrictly<T>(
  treeList: T[],
  cb: (item: T) => boolean,
  options: Required<FilterTreeOptions<T>, 'childrenKey'>,
): T[] {
  return treeList.reduce((nodes, item) => {
    if (cb(item)) {
      nodes.push(item);
    }
    if ((item as any)[options.childrenKey]) {
      return nodes.concat(filterTreeListStrictly((item as any)[options.childrenKey], cb, options));
    }
    return nodes;
  }, [] as T[]);
}

function filterTreeListAll<T>(
  treeList: T[],
  cb: (item: T) => boolean,
  options: Required<FilterTreeOptions<T>, 'childrenKey'>,
): T[] {
  const getChildren = (item: T): T[] | undefined => (item as any)[options.childrenKey];
  return treeList.filter((item) => {
    const getNewChildren = () => filterTreeListAll(getChildren(item) || [], cb, options);
    if (cb(item)) {
      if (options?.filterChildren) {
        (item as any)[options.childrenKey] = getNewChildren();
      }
      return true;
    } else {
      const children = getNewChildren();
      return children.length > 0 && ((item as any)[options.childrenKey] = children);
    }
  });
}

function filterTreeListOnlyLeaf<T>(
  treeList: T[],
  cb: (item: T) => boolean,
  options: Required<FilterTreeOptions<T>, 'childrenKey'>,
): T[] {
  const isLeaf = options?.isLeaf || ((node: T) => !(node as any)[options.childrenKey]);
  return treeList.filter((node) => {
    if (isLeaf(node)) {
      return cb(node);
    } else if ((node as any)[options.childrenKey]) {
      const children = filterTreeListOnlyLeaf((node as any)[options.childrenKey], cb, options);
      if (children && children.length > 0) {
        (node as any)[options.childrenKey] = children;
        return true;
      }
      return false;
    } else {
      return false;
    }
  });
}

export function filterTreeList<T>(
  _treeList: T[],
  cb: (item: T) => boolean,
  _options?: FilterTreeOptions<T>,
): T[] {
  const treeList = cloneDeep(_treeList);
  const options = { ..._options, ...resolveOptions(_options) };
  if (options?.onlyLeaf) {
    return filterTreeListOnlyLeaf(treeList, cb, options);
  } else if (options?.filterStrictly) {
    return filterTreeListStrictly(treeList, cb, options);
  } else {
    return filterTreeListAll(treeList, cb, options);
  }
}

export function mapTreeList<T, R>(
  treeList: T[],
  cb: (node: T, options: MapTreeOptions<T, R>) => R,
  options?: TraverseTreeOptions,
): R[] {
  return treeList?.map?.((node) => mapTree<T, R>(node, cb, options)!);
}
