import { isUrl } from "./format";
import { globalConfig } from "./global-config";

export const getResourceUrl = (originalUrl: string) => {
    // 如果本身就是一个url,或者前后端同源,直接返回
    if (isUrl(originalUrl) || !globalConfig.domain) {
        return originalUrl;
    }

    // 前后端不同源
    return `http${globalConfig.https ? 's' : ''}://${globalConfig.domain}${originalUrl}`;
}