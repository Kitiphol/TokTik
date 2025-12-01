// /app/api/videos/route.ts
import { NextRequest, NextResponse } from "next/server";
import { S3Client, PutObjectCommand } from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";

const s3Client = new S3Client({
  region: "auto",
  endpoint: `https://${process.env.R2_ACCOUNT_ID}.r2.cloudflarestorage.com`,
  credentials: {
    accessKeyId: process.env.R2_ACCESS_KEY_ID!,
    secretAccessKey: process.env.R2_SECRET_ACCESS_KEY!,
  },
});

export async function POST(req: NextRequest) {
  try {
    const { filename, title } = await req.json();

    if (!filename || !title) {
      return NextResponse.json({ error: "filename and title are required" }, { status: 400 });
    }

    const bucketName = process.env.R2_BUCKET_NAME!;
    const objectKey = `videos/${Date.now()}-${filename}`;

    const command = new PutObjectCommand({
      Bucket: bucketName,
      Key: objectKey,
      ContentType: "video/mp4",
    });

    const presignedURL = await getSignedUrl(s3Client, command, { expiresIn: 3600 });

    // Optionally store metadata in DB here

    return NextResponse.json({ presignedURL, objectKey });
  } catch (error) {
    console.error(error);
    return NextResponse.json({ error: "Failed to generate presigned URL" }, { status: 500 });
  }
}
