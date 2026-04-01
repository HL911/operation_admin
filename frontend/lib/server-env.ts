/**
 * RequiredServerEnvOptions 描述读取必填服务端环境变量时的补充提示。
 */
export interface RequiredServerEnvOptions {
  /** example 表示建议展示给开发者参考的配置示例值。 */
  example?: string;
  /** location 表示推荐开发者放置配置的文件位置。 */
  location?: string;
}

/**
 * getRequiredServerEnv 负责读取必填的服务端环境变量，并在缺失时抛出可定位的中文错误。
 */
export function getRequiredServerEnv(
  key: string,
  options: RequiredServerEnvOptions = {},
): string {
  // configuredValue 保存当前环境变量中读取到的原始值。
  const configuredValue = process.env[key]?.trim();

  if (configuredValue) {
    return configuredValue;
  }

  // locationText 用于在错误消息中提示推荐的配置文件位置。
  const locationText = options.location ? `，请在 ${options.location} 中设置` : "，请先完成配置";
  // exampleText 用于在错误消息中展示推荐的配置示例。
  const exampleText = options.example ? `，例如 ${key}=${options.example}` : "";

  throw new Error(`缺少服务端配置 ${key}${locationText}${exampleText}。`);
}
