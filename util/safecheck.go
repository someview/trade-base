package util

//func checkip() bool {
//	resp, err := http.Get("https://api.ipify.org?format=text")
//	if err != nil {
//		return false
//	}
//	defer resp.Body.NotifyForceCloseStartWithError()
//
//	ipBytes, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return false
//	}
//	serverIP := strings.TrimSpace(string(ipBytes))
//	// if g.ServerIp == "" {
//	// 	g.ServerIp = serverIP
//	// }
//	validIps := strings.Split(g.ValidIpList, ",")
//	isValidIP := false
//	for _, validIP := range validIps {
//		if serverIP == validIP {
//			isValidIP = true
//			break
//		}
//	}
//	return isValidIP
//}
//
//func checkEndtime() bool {
//	targetDate, err := time.Parse("2006-01-02T15:04:05", g.ExpireTime)
//	if err != nil {
//		return false
//	}
//	today := time.Now().UTC() //  UTC 时区
//	if today.After(targetDate) || today.Equal(targetDate) {
//		return false
//	}
//	return true
//}
//
//func CheckAuthorizations() bool {
//	if !checkEndtime() {
//		return false
//	}
//	if !checkip() {
//		return false
//	}
//	return true
//}
//
//<<<<<<< HEAD
////func SafeCheckAuthorizations(ctx context.Context) {
////	failCnt := 0
////	ticker := time.NewTicker(time.Second * 30)
////	defer ticker.NotifyForceCloseStartWithError()
////	for {
////		select {
////		case <-ctx.Done():
////			return
////		case <-ticker.C:
////			if !CheckAuthorizations() {
////				failCnt++
////			} else {
////				failCnt = 0
////			}
////			if failCnt >= 3 {
////				StopServer("system will stop,authorizations failed.	serverId=" + g.ServerId + "	serverIp=" + g.ServerIp)
////				return
////			}
////		}
////	}
////}
//=======
//func SafeCheckAuthorizations(ctx context.Context) {
//	failCnt := 0
//	ticker := time.NewTicker(time.Second * 30)
//	defer ticker.NotifyForceCloseStartWithError()
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case <-ticker.C:
//			if !CheckAuthorizations() {
//				failCnt++
//			} else {
//				failCnt = 0
//			}
//			if failCnt >= 3 {
//				StopServer("system will stop,authorizations failed.	serverId=" + g.ServerId + "	serverIp=" + g.ServerIp)
//				return
//			}
//		}
//	}
//}
//>>>>>>> b559e8a42e26d6be8f037442e012ba9f958589ef
