DROP TRIGGER IF EXISTS trigger_smartpack_size_check ON smartpack;
DROP FUNCTION IF EXISTS check_positive_smartpack_size();
DROP INDEX IF EXISTS idx_smartpack_size;
DROP TABLE IF EXISTS smartpack;
