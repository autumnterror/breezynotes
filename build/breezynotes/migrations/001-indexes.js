const dbName = process.env.MONGO_INITDB_DATABASE || "blocknotedb";
const dbRef = db.getSiblingDB(dbName);

print("Applying Mongo indexes...");

dbRef.notes.createIndex(
  { author: 1, updated_at: -1 },
  { name: "idx_notes_author_updatedAt" },
);

dbRef.notes.createIndex(
  { editors: 1, updated_at: -1 },
  { name: "idx_notes_editors_updatedAt" },
);

dbRef.notes.createIndex(
  { readers: 1, updated_at: -1 },
  { name: "idx_notes_readers_updatedAt" },
);

dbRef.notes.createIndex(
  { is_public: 1, updated_at: -1 },
  { name: "idx_notes_public_updatedAt" },
);

dbRef.notes.createIndex(
  { is_blog: 1, updated_at: -1 },
  { name: "idx_notes_blog_updatedAt" },
);

dbRef.blocks.createIndex(
  { note_id: 1, created_at: 1 },
  { name: "idx_blocks_note_createdAt" },
);

dbRef.blocks.createIndex(
  { is_used: 1, updated_at: 1 },
  { name: "idx_blocks_isUsed_updatedAt" },
);

dbRef.tags.createIndex(
  { user_id: 1, title: 1 },
  {
    name: "uniq_tags_user_title",
    unique: true,
  },
);

dbRef.usertags.createIndex(
  { note_id: 1, tag: 1 },
  {
    name: "uniq_usertags_note_tag",
    unique: true,
  },
);

dbRef.usertags.createIndex({ tag: 1 }, { name: "idx_usertags_tag" });

dbRef.migrations.updateOne(
  { _id: "001-indexes" },
  { $setOnInsert: { appliedAt: new Date() } },
  { upsert: true },
);

print("Mongo indexes applied successfully ✅");
