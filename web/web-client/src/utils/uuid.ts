export const createUUID = () =>{
  const tmpUrl = URL.createObjectURL(new Blob());
  const uuid = tmpUrl.split('/').pop() as string;
  URL.revokeObjectURL(tmpUrl);
  return uuid;
}