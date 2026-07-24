package core

// Source ইন্টারফেসটি ডিফাইন করা হলো। 
// ভবিষ্যতে আমরা যতগুলো নতুন API যোগ করব, সবাইকে এই ইন্টারফেস মেনে চলতে হবে।
type Source interface {
	Name() string                           // সোর্সের নাম রিটার্ন করবে
	Run(domain string) ([]string, error)    // ডোমেইন নিয়ে সাবডোমেইনের লিস্ট রিটার্ন করবে
}
