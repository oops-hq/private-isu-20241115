alter table users add index idx_account_name_del_flg (account_name, del_flg);
alter table comments add index idx_comments_post_id_created_at (post_id, created_at);
alter table comments add index idx_comments_user_id (user_id);
alter table posts add index idx_posts_created_at (created_at);
alter table posts add index idx_posts_user_id_created_at (user_id, created_at);