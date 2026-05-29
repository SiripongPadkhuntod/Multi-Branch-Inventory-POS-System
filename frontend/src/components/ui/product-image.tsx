import { ImageIcon } from "lucide-react";

type ProductImageProps = {
  src?: string;
  name: string;
  className?: string;
};

export function ProductImage({ src, name, className = "h-14 w-14" }: ProductImageProps) {
  if (src) {
    return (
      <img
        src={src}
        alt={name}
        className={`${className} shrink-0 rounded-md border border-line bg-white object-cover`}
        loading="lazy"
      />
    );
  }

  const initials = name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");

  return (
    <div className={`${className} grid shrink-0 place-items-center rounded-md border border-line bg-brandSoft text-brand`}>
      {initials ? <span className="text-sm font-bold">{initials}</span> : <ImageIcon className="h-5 w-5" />}
    </div>
  );
}
