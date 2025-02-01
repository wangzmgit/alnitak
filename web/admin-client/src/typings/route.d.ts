interface BaseRouteType {
  name: string;
  path: string;
  // component: string;
  // children?: RouteType[]
}

interface RouteType extends BaseRouteType {
  component: string;
  children?: RouteType[]
}

interface AsyncRouteType extends BaseRouteType {
  component?: () => Promise<any>;
  children?: AsyncRouteType[],
}